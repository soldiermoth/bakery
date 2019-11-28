package filters

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/cbsinteractive/bakery/config"
	"github.com/cbsinteractive/bakery/parsers"
	"github.com/grafov/m3u8"
)

// HLSFilter implements the Filter interface for HLS
// manifests
type HLSFilter struct {
	manifestURL     string
	manifestContent string
	config          config.Config
}

// NewHLSFilter is the HLS filter constructor
func NewHLSFilter(manifestURL, manifestContent string, c config.Config) *HLSFilter {
	return &HLSFilter{
		manifestURL:     manifestURL,
		manifestContent: manifestContent,
		config:          c,
	}
}

// FilterManifest will be responsible for filtering the manifest
// according  to the MediaFilters
func (h *HLSFilter) FilterManifest(filters *parsers.MediaFilters) (string, error) {
	m, manifestType, err := m3u8.DecodeFrom(strings.NewReader(h.manifestContent), true)
	if err != nil {
		return "", err
	}
	if manifestType != m3u8.MASTER {
		return "", errors.New("manifest type is wrong")
	}

	// convert into the master playlist type
	manifest := m.(*m3u8.MasterPlaylist)
	filteredManifest := m3u8.NewMasterPlaylist()

	for _, v := range manifest.Variants {
		// transform media manifests and alternatives (EXT-X-MEDIA)
		// into absolute urls
		absoluteURL, _ := filepath.Split(h.manifestURL)
		if len(v.VariantParams.Alternatives) > 0 {
			for _, a := range v.VariantParams.Alternatives {
				a.URI = absoluteURL + a.URI
			}
		}
		filteredManifest.Append(absoluteURL+v.URI, v.Chunklist, v.VariantParams)
	}
	return filteredManifest.String(), nil
}
