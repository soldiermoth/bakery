package parser

import (
	"encoding/json"
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
				Videos: []VideoType{videoHDR10},
			},
		},
		{
			"two video types",
			"/v(hdr10,hevc)/",
			MediaFilters{
				Videos: []VideoType{videoHDR10, videoHEVC},
			},
		},
		{
			"two video types and two audio types",
			"/v(hdr10,hevc)/a(aac,noAd)/",
			MediaFilters{
				Videos: []VideoType{videoHDR10, videoHEVC},
				Audios: []AudioType{audioAAC, audioNoAudioDescription},
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
			"detect protocol hls for urls with .m3u8 extension",
			"url/here/with/master.m3u8",
			MediaFilters{
				Protocol: protocolHLS,
			},
		},
		{
			"detect protocol dash for urls with .mpd extension",
			"url/here/with/manifest.mpd",
			MediaFilters{
				Protocol: protocolDASH,
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
