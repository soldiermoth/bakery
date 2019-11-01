package parser

import (
	"regexp"
	"strings"
)

// VideoType is the video codec we need in a given playlist
type VideoType string

// AudioType is the audio codec we need in a given playlist
type AudioType string

// AudioLanguage is the audio language we need in a given playlist
type AudioLanguage string

// CaptionLanguage is the audio language we need in a given playlist
type CaptionLanguage string

const (
	// Video Types

	VideoHDR10       VideoType = "hdr10"
	VideoDolbyVision VideoType = "dovi"
	VideoHEVC        VideoType = "hevc"
	VideoH264        VideoType = "avc"

	// Audio Types

	AudioAAC                AudioType = "aac"
	AudioNoAudioDescription AudioType = "noAd"

	// Audio Languages

	AudioPTBR AudioLanguage = "pt-br"
	AudioES   AudioLanguage = "es"
	AudioEN   AudioLanguage = "en"

	// Captions Languages

	CaptionPTBR CaptionLanguage = "pt-br"
	CaptionES   CaptionLanguage = "es"
	CaptionEN   CaptionLanguage = "en"
)

// MediaFilters is a struct that carry all the information passed via url
type MediaFilters struct {
	Videos           []VideoType       `json:"Videos,omitempty"`
	Audios           []AudioType       `json:"Audios,omitempty"`
	AudioLanguages   []AudioLanguage   `json:"AudioLanguages,omitempty"`
	CaptionLanguages []CaptionLanguage `json:"CaptionLanguages,omitempty"`
	MaxBitrate       int               `json:"MaxBitrate,omitempty"`
	MinBitrate       int               `json:"MaxBitrate,omitempty"`
}

func parse(filters string) (*MediaFilters, error) {
	mf := new(MediaFilters)
	parts := strings.Split(filters, "/")
	re := regexp.MustCompile(`(.*)\((.*)\)`)

	for _, part := range parts {
		subparts := re.FindStringSubmatch(part)
		// FindStringSubmatch should return a slice with
		// the full string, the key and values (3 elements)
		if len(subparts) != 3 {
			continue
		}

		switch key := subparts[1]; key {
		case "video":
			values := strings.Split(subparts[2], ",")
			for _, videoType := range values {
				mf.Videos = append(mf.Videos, VideoType(videoType))
			}
		}

	}

	return mf, nil
}
