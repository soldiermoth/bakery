package parsers

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

func TestURLParseUrl(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		expectedFilters      MediaFilters
		expectedManifestPath string
	}{
		{
			"one video type",
			"/v(hdr10)/",
			MediaFilters{
				Videos:     []VideoType{"hev1.2", "hvc1.2"},
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
				VideoSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
					Codecs:     []Codec{codecHDR10},
				},
				AudioSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
			},
			"/",
		},
		{
			"two video types",
			"/v(hdr10,hevc)/",
			MediaFilters{
				Videos:     []VideoType{"hev1.2", "hvc1.2", videoHEVC},
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
				VideoSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
					Codecs:     []Codec{codecHDR10, codecHEVC},
				},
				AudioSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
			},
			"/",
		},
		{
			"two video types and two audio types",
			"/v(hdr10,hevc)/a(aac,noAd)/",
			MediaFilters{
				Videos:     []VideoType{"hev1.2", "hvc1.2", videoHEVC},
				Audios:     []AudioType{audioAAC, audioNoAudioDescription},
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
				VideoSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
					Codecs:     []Codec{codecHDR10, codecHEVC},
				},
				AudioSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
					Codecs:     []Codec{codecAAC, codecNoAudioDescription},
				},
			},
			"/",
		},
		{
			"videos, audio, captions and bitrate range",
			"/v(hdr10,hevc)/a(aac)/al(pt-BR,en)/c(en)/b(100,4000)/",
			MediaFilters{
				Videos:           []VideoType{"hev1.2", "hvc1.2", videoHEVC},
				Audios:           []AudioType{audioAAC},
				AudioLanguages:   []AudioLanguage{audioLangPTBR, audioLangEN},
				CaptionLanguages: []CaptionLanguage{captionEN},
				MaxBitrate:       4000,
				MinBitrate:       100,
				VideoSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
					Codecs:     []Codec{codecHDR10, codecHEVC},
				},
				AudioSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
					Codecs:     []Codec{codecAAC},
				},
			},
			"/",
		},
		{
			"bitrate range with minimum bitrate only",
			"/b(100,)/",
			MediaFilters{
				MaxBitrate: math.MaxInt32,
				MinBitrate: 100,
				VideoSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
				AudioSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
			},
			"/",
		},
		{
			"bitrate range with maximum bitrate only",
			"/b(,3000)/",
			MediaFilters{
				MaxBitrate: 3000,
				MinBitrate: 0,
				VideoSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
				AudioSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
			},
			"/",
		},
		{
			"detect protocol hls for urls with .m3u8 extension",
			"/path/here/with/master.m3u8",
			MediaFilters{
				Protocol:   ProtocolHLS,
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
				VideoSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
				AudioSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
			},
			"/path/here/with/master.m3u8",
		},
		{
			"detect protocol dash for urls with .mpd extension",
			"/path/here/with/manifest.mpd",
			MediaFilters{
				Protocol:   ProtocolDASH,
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
				VideoSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
				AudioSubFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
			},
			"/path/here/with/manifest.mpd",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			masterManifestPath, output, err := URLParse(test.input)
			if err != nil {
				t.Fatal(err)
			}

			jsonOutput, err := json.Marshal(output)
			if err != nil {
				t.Fatal(err)
			}

			jsonExpected, err := json.Marshal(test.expectedFilters)
			if err != nil {
				t.Fatal(err)
			}

			if test.expectedManifestPath != masterManifestPath {
				t.Errorf("wrong master manifest generated.\nwant %#v\n\ngot %#v", test.expectedManifestPath, masterManifestPath)
			}

			if !reflect.DeepEqual(jsonOutput, jsonExpected) {
				t.Errorf("wrong struct generated.\nwant %#v\ngot %#v", test.expectedFilters, output)
			}
		})
	}
}
