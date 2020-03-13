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
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
				VideoFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
					Codecs:     []Codec{"hev1.2", "hvc1.2"},
				},
				AudioFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
			},
			"/",
		},
		{
			"two video types",
			"/v(hdr10,hevc)/",
			MediaFilters{
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
				VideoFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
					Codecs:     []Codec{"hev1.2", "hvc1.2", codecHEVC},
				},
				AudioFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
			},
			"/",
		},
		{
			"two video types and two audio types",
			"/v(hdr10,hevc)/a(aac,noAd)/",
			MediaFilters{
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
				VideoFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
					Codecs:     []Codec{"hev1.2", "hvc1.2", codecHEVC},
				},
				AudioFilters: Subfilters{
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
				AudioLanguages:   []AudioLanguage{audioLangPTBR, audioLangEN},
				CaptionLanguages: []CaptionLanguage{captionEN},
				MaxBitrate:       4000,
				MinBitrate:       100,
				VideoFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
					Codecs:     []Codec{"hev1.2", "hvc1.2", codecHEVC},
				},
				AudioFilters: Subfilters{
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
				VideoFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
				AudioFilters: Subfilters{
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
				VideoFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
				AudioFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
			},
			"/",
		},
		{
			"nested audio and video bitrate filters",
			"/a(b(100,))/v(b(,5000))/",
			MediaFilters{
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
				VideoFilters: Subfilters{
					MaxBitrate: 5000,
				},
				AudioFilters: Subfilters{
					MinBitrate: 100,
					MaxBitrate: math.MaxInt32,
				},
			},
			"/",
		},
		{
			"nested codec and bitrate filters in audio",
			"/a(b(100,200),codec(ac-3,aac))/",
			MediaFilters{
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
				VideoFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
				AudioFilters: Subfilters{
					MinBitrate: 100,
					MaxBitrate: 200,
					Codecs:     []Codec{codecAC3, codecAAC},
				},
			},
			"/",
		},
		{
			"nested codec and bitrate filters in video, plus overall bitrate filters",
			"/v(codec(avc,hdr10),b(1000,2000))/",
			MediaFilters{
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
				VideoFilters: Subfilters{
					MaxBitrate: 2000,
					MinBitrate: 1000,
					Codecs:     []Codec{codecH264, "hev1.2", "hvc1.2"},
				},
				AudioFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
			},
			"/",
		},
		{
			"nested bitrate and old format of codec filter",
			"/a(mp4a,ac-3,b(0,10))/",
			MediaFilters{
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
				VideoFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
				AudioFilters: Subfilters{
					MaxBitrate: 10,
					Codecs:     []Codec{"mp4a", codecAC3},
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
				VideoFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
				AudioFilters: Subfilters{
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
				VideoFilters: Subfilters{
					MaxBitrate: math.MaxInt32,
				},
				AudioFilters: Subfilters{
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
