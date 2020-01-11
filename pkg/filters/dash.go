package filters

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/cbsinteractive/bakery/pkg/config"
	"github.com/cbsinteractive/bakery/pkg/parsers"
	"github.com/zencoder/go-dash/mpd"
)

const adaptationSetTypeText = "text"

// DASHFilter implements the Filter interface for DASH manifests
type DASHFilter struct {
	manifestURL     string
	manifestContent string
	config          config.Config
}

// NewDASHFilter is the DASH filter constructor
func NewDASHFilter(manifestURL, manifestContent string, c config.Config) *DASHFilter {
	return &DASHFilter{
		manifestURL:     manifestURL,
		manifestContent: manifestContent,
		config:          c,
	}
}

// FilterManifest will be responsible for filtering the manifest according  to the MediaFilters
func (d *DASHFilter) FilterManifest(filters *parsers.MediaFilters) (string, error) {
	manifest, err := mpd.ReadFromString(d.manifestContent)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(d.manifestURL)
	if err != nil {
		return "", fmt.Errorf("parsing manifest url: %w", err)
	}

	baseURLWithPath := func(p string) string {
		var sb strings.Builder
		sb.WriteString(u.Scheme)
		sb.WriteString("://")
		sb.WriteString(u.Host)
		sb.WriteString(p)
		sb.WriteString("/")
		return sb.String()
	}

	if manifest.BaseURL == "" {
		manifest.BaseURL = baseURLWithPath(path.Dir(u.Path))
	} else if !strings.HasPrefix(manifest.BaseURL, "http") {
		manifest.BaseURL = baseURLWithPath(path.Join(path.Dir(u.Path), manifest.BaseURL))
	}

	if filters.CaptionTypes != nil {
		d.filterCaptionTypes(filters, manifest)
	}

	return manifest.WriteToString()
}

func (d *DASHFilter) filterCaptionTypes(filters *parsers.MediaFilters, manifest *mpd.MPD) {
	supportedTypes := map[parsers.CaptionType]struct{}{}

	for _, captionType := range filters.CaptionTypes {
		supportedTypes[captionType] = struct{}{}
	}

	for _, period := range manifest.Periods {
		for _, as := range period.AdaptationSets {
			if as.ContentType == nil {
				continue
			}

			if *as.ContentType == adaptationSetTypeText {
				var filteredReps []*mpd.Representation

				for _, r := range as.Representations {
					if r.Codecs == nil {
						filteredReps = append(filteredReps, r)
						continue
					}

					if _, supported := supportedTypes[parsers.CaptionType(*r.Codecs)]; supported {
						filteredReps = append(filteredReps, r)
					}
				}

				as.Representations = filteredReps
			}
		}
	}
}
