package filters

import "github.com/cbsinteractive/bakery/parsers"

// Filter is an interface for HLS and DASH filters
type Filter interface {
	FilterManifest(filters *parsers.MediaFilters) (string, error)
}
