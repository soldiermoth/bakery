package parsers

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

func TestURLParseUrl(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected MediaFilters
	}{
		{
			"one video type",
			"/v(hdr10)/",
			MediaFilters{
				Videos:     []VideoType{videoHDR10},
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
			},
		},
		{
			"two video types",
			"/v(hdr10,hevc)/",
			MediaFilters{
				Videos:     []VideoType{videoHDR10, videoHEVC},
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
			},
		},
		{
			"two video types and two audio types",
			"/v(hdr10,hevc)/a(aac,noAd)/",
			MediaFilters{
				Videos:     []VideoType{videoHDR10, videoHEVC},
				Audios:     []AudioType{audioAAC, audioNoAudioDescription},
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
			},
		},
		{
			"videos, audio, captions and bitrate range",
			"/v(hdr10,hevc)/a(aac)/al(pt-br,en)/c(en)/b(100,4000)/",
			MediaFilters{
				Videos:           []VideoType{videoHDR10, videoHEVC},
				Audios:           []AudioType{audioAAC},
				AudioLanguages:   []AudioLanguage{audioLangPTBR, audioLangEN},
				CaptionLanguages: []CaptionLanguage{captionEN},
				MaxBitrate:       4000,
				MinBitrate:       100,
			},
		},
		{
			"bitrate range with minimum bitrate only",
			"/b(100,)/",
			MediaFilters{
				MaxBitrate: math.MaxInt32,
				MinBitrate: 100,
			},
		},
		{
			"bitrate range with maximum bitrate only",
			"/b(,3000)/",
			MediaFilters{
				MaxBitrate: 3000,
				MinBitrate: 0,
			},
		},
		{
			"detect protocol hls for urls with .m3u8 extension",
			"url/here/with/master.m3u8",
			MediaFilters{
				Protocol:   ProtocolHLS,
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
			},
		},
		{
			"detect protocol dash for urls with .mpd extension",
			"url/here/with/manifest.mpd",
			MediaFilters{
				Protocol:   ProtocolDASH,
				MaxBitrate: math.MaxInt32,
				MinBitrate: 0,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			output, err := URLParse(test.input)
			if err != nil {
				t.Fatal(err)
			}

			jsonOutput, err := json.Marshal(output)
			if err != nil {
				t.Fatal(err)
			}

			jsonExpected, err := json.Marshal(test.expected)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(jsonOutput, jsonExpected) {
				t.Errorf("wrong struct generated.\nwant %#v\ngot %#v", test.expected, output)
			}
		})
	}
}
