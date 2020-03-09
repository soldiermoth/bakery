package filters

import (
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/cbsinteractive/bakery/pkg/config"
	"github.com/cbsinteractive/bakery/pkg/parsers"
	"github.com/zencoder/go-dash/mpd"
)

type execFilter func(filters *parsers.MediaFilters, manifest *mpd.MPD)

// DASHFilter implements the Filter interface for DASH manifests
type DASHFilter struct {
	manifestURL     string
	manifestContent string
	config          config.Config
	filters         []execFilter
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

	for _, filter := range d.getFilters(filters) {
		filter(filters, manifest)
	}

	return manifest.WriteToString()
}

func (d *DASHFilter) getFilters(filters *parsers.MediaFilters) []execFilter {
	filterList := []execFilter{}
	if filters.FilterStreamTypes != nil && len(filters.FilterStreamTypes) > 0 {
		filterList = append(filterList, d.filterAdaptationSetType)
	}

	if filters.DefinesBitrateFilter() {
		filterList = append(filterList, d.filterBandwidth)
	}

	if filters.Videos != nil {
		filterList = append(filterList, d.filterVideoTypes)
	}

	if filters.Audios != nil {
		filterList = append(filterList, d.filterAudioTypes)
	}

	if filters.CaptionTypes != nil {
		filterList = append(filterList, d.filterCaptionTypes)
	}

	if filters.Role != "" {
		filterList = append(filterList, d.updateRoleDescription)
	}

	return filterList
}

func (d *DASHFilter) filterVideoTypes(filters *parsers.MediaFilters, manifest *mpd.MPD) {
	supportedVideoTypes := map[string]struct{}{}
	for _, videoType := range filters.Videos {
		supportedVideoTypes[string(videoType)] = struct{}{}
	}

	filterContentType(videoContentType, supportedVideoTypes, manifest)
}

func (d *DASHFilter) filterAudioTypes(filters *parsers.MediaFilters, manifest *mpd.MPD) {
	supportedAudioTypes := map[string]struct{}{}
	for _, audioType := range filters.Audios {
		supportedAudioTypes[string(audioType)] = struct{}{}
	}

	filterContentType(audioContentType, supportedAudioTypes, manifest)
}

func (d *DASHFilter) filterCaptionTypes(filters *parsers.MediaFilters, manifest *mpd.MPD) {
	supportedCaptionTypes := map[string]struct{}{}
	for _, captionType := range filters.CaptionTypes {
		supportedCaptionTypes[string(captionType)] = struct{}{}
	}

	filterContentType(captionContentType, supportedCaptionTypes, manifest)
}

func filterContentType(filter ContentType, supportedContentTypes map[string]struct{}, manifest *mpd.MPD) {
	for _, period := range manifest.Periods {
		var filteredAdaptationSets []*mpd.AdaptationSet
		for _, as := range period.AdaptationSets {
			if as.ContentType != nil && *as.ContentType == string(filter) {
				var filteredReps []*mpd.Representation
				for _, r := range as.Representations {
					if r.Codecs == nil {
						filteredReps = append(filteredReps, r)
						continue
					}

					if matchCodec(*r.Codecs, filter, supportedContentTypes) {
						continue
					}

					filteredReps = append(filteredReps, r)
				}
				as.Representations = filteredReps
			}

			if len(as.Representations) != 0 {
				filteredAdaptationSets = append(filteredAdaptationSets, as)
			}
		}

		for i, as := range filteredAdaptationSets {
			as.ID = strptr(strconv.Itoa(i))
		}
		period.AdaptationSets = filteredAdaptationSets
	}
}

func (d *DASHFilter) filterAdaptationSetType(filters *parsers.MediaFilters, manifest *mpd.MPD) {
	filteredAdaptationSetTypes := map[parsers.StreamType]struct{}{}
	for _, streamType := range filters.FilterStreamTypes {
		filteredAdaptationSetTypes[streamType] = struct{}{}
	}

	periodIndex := 0
	var filteredPeriods []*mpd.Period
	for _, period := range manifest.Periods {
		var filteredAdaptationSets []*mpd.AdaptationSet
		asIndex := 0
		for _, as := range period.AdaptationSets {
			if as.ContentType != nil {
				if _, filtered := filteredAdaptationSetTypes[parsers.StreamType(*as.ContentType)]; filtered {
					continue
				}
			}

			as.ID = strptr(strconv.Itoa(asIndex))
			asIndex++

			filteredAdaptationSets = append(filteredAdaptationSets, as)
		}

		if len(filteredAdaptationSets) == 0 {
			continue
		}

		period.AdaptationSets = filteredAdaptationSets
		period.ID = strconv.Itoa(periodIndex)
		periodIndex++

		filteredPeriods = append(filteredPeriods, period)
	}

	manifest.Periods = filteredPeriods
}

func matchCodec(codec string, ct ContentType, supportedCodecs map[string]struct{}) bool {
	//the key in supportedCodecs for captionContentType is equivalent to codec
	//advertised in manifest. we can avoid iterating through each key
	if ct == captionContentType {
		_, found := supportedCodecs[codec]
		return found
	}

	for key := range supportedCodecs {
		if ValidCodecs(codec, CodecFilterID(key)) {
			return true
		}
	}

	return false
}

func (d *DASHFilter) filterBandwidth(filters *parsers.MediaFilters, manifest *mpd.MPD) {
	for _, period := range manifest.Periods {
		var filteredAdaptationSets []*mpd.AdaptationSet

		for _, as := range period.AdaptationSets {
			var filteredRepresentations []*mpd.Representation

			for _, r := range as.Representations {
				if r.Bandwidth == nil {
					continue
				}
				if *r.Bandwidth <= int64(filters.MaxBitrate) && *r.Bandwidth >= int64(filters.MinBitrate) {
					filteredRepresentations = append(filteredRepresentations, r)
				}
			}
			as.Representations = filteredRepresentations
			if len(as.Representations) != 0 {
				filteredAdaptationSets = append(filteredAdaptationSets, as)
			}
		}

		period.AdaptationSets = filteredAdaptationSets

		// Recalculate AdaptationSet id numbers
		for index, as := range period.AdaptationSets {
			as.ID = strptr(strconv.Itoa(index))
		}
	}
}

//
func (d *DASHFilter) updateRoleDescription(filters *parsers.MediaFilters, manifest *mpd.MPD) {
	for _, period := range manifest.Periods {
		for _, as := range period.AdaptationSets {
			for i, accessibility := range as.AccessibilityElems {
				if *accessibility.SchemeIdUri == "urn:tva:metadata:cs:AudioPurposeCS:2007" {
					as.Roles[i].Value = strptr(filters.Role)
				}
			}
		}
	}
}

func strptr(s string) *string {
	return &s
}
