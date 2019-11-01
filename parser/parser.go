package parser

type VideoType string
type AudioType string
type AudioLanguage string
type CaptionLanguage string

const (
	// Video Types
	VideoHDR10       VideoType = "hdr"
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

	// Caption Languages
	CaptionPTBR CaptionLanguage = "pt-br"
	CaptionES   CaptionLanguage = "es"
	CaptionEN   CaptionLanguage = "en"
)

type MediaFilters struct {
	Videos           []VideoType       `json:"Videos,omitempty"`
	Audios           []AudioType       `json:"Audios,omitempty"`
	AudioLanguages   []AudioLanguage   `json:"AudioLanguages,omitempty"`
	CaptionLanguages []CaptionLanguage `json:"CaptionLanguages,omitempty"`
	MaxBitrate       int               `json:"MaxBitrate,omitempty"`
	MinBitrate       int               `json:"MaxBitrate,omitempty"`
}

func parse(filters string) (*MediaFilters, error) {
	return new(MediaFilters), nil
}
