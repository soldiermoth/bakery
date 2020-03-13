package filters

import (
	"fmt"

	"github.com/cbsinteractive/bakery/pkg/parsers"
	"github.com/grafov/m3u8"
)

// FilterManifest will be responsible for filtering the manifest
// according  to the MediaFilters
func (h *HLSFilter) filterRenditionManifest(filters *parsers.MediaFilters, m *m3u8.MediaPlaylist) (string, error) {
	fmt.Println(m)
	fmt.Println(filters)

	return "", nil
}
