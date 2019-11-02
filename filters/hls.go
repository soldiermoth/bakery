package filters

import (
	"github.com/cbsinteractive/bakery/config"
	"github.com/cbsinteractive/bakery/parsers"
	"github.com/quangngotan95/go-m3u8/m3u8"
)

// HLSFilter implements the Filter interface for HLS
// manifests
type HLSFilter struct {
	manifestContent string
	config          config.Config
}

// NewHLSFilter is the HLS filter constructor
func NewHLSFilter(manifestContent string, c config.Config) *HLSFilter {
	return &HLSFilter{
		manifestContent: manifestContent,
		config:          c,
	}
}

// FilterManifest will be responsible for filtering the manifest
// according  to the MediaFilters
func (h *HLSFilter) FilterManifest(filters *parsers.MediaFilters) (string, error) {
	manifest, err := m3u8.ReadString(h.manifestContent)
	if err != nil {
		return "", err
	}

	return manifest.String(), nil
}
