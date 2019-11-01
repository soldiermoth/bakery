package parser

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseUrl(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected MediaFilters
	}{
		{
			"one video type",
			"/video(hdr10)/",
			MediaFilters{
				Videos: []VideoType{VideoHDR10},
			},
		},
	}

	//tests := map[string]MediaFilters{
	//filters := "/video(hdr10,hevc)/audio(pt-br,en)/tracks(en)/bitrates(100,4000)/"
	//expected := MediaFilters{
	//Videos:           []VideoType{VideoHDR10, VideoHEVC},
	//AudioLanguages:   []AudioLanguage{AudioPTBR, AudioEN},
	//CaptionLanguages: []CaptionLanguage{CaptionEN},
	//MaxBitrate:       4000,
	//MinBitrate:       100,
	//}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			output, err := parse(test.input)
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
