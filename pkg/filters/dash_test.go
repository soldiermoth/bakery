package filters

import (
	"testing"

	"github.com/cbsinteractive/bakery/pkg/config"
	"github.com/cbsinteractive/bakery/pkg/parsers"
	"github.com/google/go-cmp/cmp"
)

func TestDASHFilter_FilterManifest_captionTypes(t *testing.T) {
	manifestWithWVTTAndSTPPCaptions := `<?xml version="1.0" encoding="UTF-8"?>
<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" profiles="urn:mpeg:dash:profile:isoff-on-demand:2011" type="static" mediaPresentationDuration="PT6M16S" minBufferTime="PT1.97S">
  <Period>
    <AdaptationSet id="7357" lang="en" contentType="text">
      <Representation bandwidth="256" codecs="wvtt" id="subtitle_en"></Representation>
      <Representation bandwidth="256" codecs="stpp" id="subtitle_en_ttml"></Representation>
    </AdaptationSet>
  </Period>
</MPD>
`

	manifestWithWVTTCaptions := `<?xml version="1.0" encoding="UTF-8"?>
<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" profiles="urn:mpeg:dash:profile:isoff-on-demand:2011" type="static" mediaPresentationDuration="PT6M16S" minBufferTime="PT1.97S">
  <Period>
    <AdaptationSet id="7357" lang="en" contentType="text">
      <Representation bandwidth="256" codecs="wvtt" id="subtitle_en"></Representation>
    </AdaptationSet>
  </Period>
</MPD>
`

	manifestWithSTPPCaptions := `<?xml version="1.0" encoding="UTF-8"?>
<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" profiles="urn:mpeg:dash:profile:isoff-on-demand:2011" type="static" mediaPresentationDuration="PT6M16S" minBufferTime="PT1.97S">
  <Period>
    <AdaptationSet id="7357" lang="en" contentType="text">
      <Representation bandwidth="256" codecs="stpp" id="subtitle_en_ttml"></Representation>
    </AdaptationSet>
  </Period>
</MPD>
`

	manifestWithoutCaptions := `<?xml version="1.0" encoding="UTF-8"?>
<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" profiles="urn:mpeg:dash:profile:isoff-on-demand:2011" type="static" mediaPresentationDuration="PT6M16S" minBufferTime="PT1.97S">
  <Period>
    <AdaptationSet id="7357" lang="en" contentType="text"></AdaptationSet>
  </Period>
</MPD>
`

	tests := []struct {
		name                  string
		filters               *parsers.MediaFilters
		manifestContent       string
		expectManifestContent string
		expectErr             bool
	}{
		{
			name: "when an empty caption type filter list is supplied, captions are stripped from a " +
				"manifest",
			filters:               &parsers.MediaFilters{CaptionTypes: []parsers.CaptionType{}},
			manifestContent:       manifestWithWVTTAndSTPPCaptions,
			expectManifestContent: manifestWithoutCaptions,
		},
		{
			name: "when a caption type filter is supplied with stpp only, webvtt captions are " +
				"filtered out",
			filters:               &parsers.MediaFilters{CaptionTypes: []parsers.CaptionType{"stpp"}},
			manifestContent:       manifestWithWVTTAndSTPPCaptions,
			expectManifestContent: manifestWithSTPPCaptions,
		},
		{
			name: "when a caption type filter is supplied with wvtt only, stpp captions are " +
				"filtered out",
			filters:               &parsers.MediaFilters{CaptionTypes: []parsers.CaptionType{"wvtt"}},
			manifestContent:       manifestWithWVTTAndSTPPCaptions,
			expectManifestContent: manifestWithWVTTCaptions,
		},
		{
			name:                  "when no filters are supplied, captions are not stripped from a manifest",
			filters:               &parsers.MediaFilters{},
			manifestContent:       manifestWithWVTTAndSTPPCaptions,
			expectManifestContent: manifestWithWVTTAndSTPPCaptions,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			filter := NewDASHFilter("", tt.manifestContent, config.Config{})

			manifest, err := filter.FilterManifest(tt.filters)
			if err != nil && !tt.expectErr {
				t.Errorf("FilterManifest() didnt expect an error to be returned, got: %v", err)
				return
			} else if err == nil && tt.expectErr {
				t.Error("FilterManifest() expected an error, got nil")
				return
			}

			if g, e := manifest, tt.expectManifestContent; g != e {
				t.Errorf("FilterManifest() wrong manifest returned\ngot %v\nexpected: %v\ndiff: %v", g, e,
					cmp.Diff(g, e))
			}
		})
	}
}

func strptr(str string) *string {
	return &str
}
