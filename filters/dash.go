package filters

import (
	"github.com/cbsinteractive/bakery/config"
	"github.com/cbsinteractive/bakery/parsers"
	"github.com/cbsinteractive/go-dash/mpd"
)

// DASHFilter implements the Filter interface for DASH
// manifests
type DASHFilter struct {
	manifestContent string
	config          config.Config
}

// NewDASHFilter is the DASH filter constructor
func NewDASHFilter(manifestContent string, c config.Config) *DASHFilter {
	return &DASHFilter{
		manifestContent: manifestContent,
		config:          c,
	}
}

// FilterManifest will be responsible for filtering the manifest
// according  to the MediaFilters
func (d *DASHFilter) FilterManifest(filters *parsers.MediaFilters) (string, error) {
	manifest, err := mpd.ReadFromString(d.manifestContent)
	if err != nil {
		return "", err
	}

	return manifest.WriteToString()
}
