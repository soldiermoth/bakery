package parsers

import (
	"math"
	"path"
	"regexp"
	"strconv"
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

// CaptionType is an allowed caption format for the stream
type CaptionType string

// StreamType represents one stream type (e.g. video, audio, text)
type StreamType string

type Codec string

// Protocol describe the valid protocols
type Protocol string

const (
	videoHDR10       VideoType = "hdr10"
	videoDolbyVision VideoType = "dovi"
	videoHEVC        VideoType = "hevc"
	videoH264        VideoType = "avc"

	audioAAC                AudioType = "aac"
	audioAC3                AudioType = "ac-3"
	audioEnhacedAC3         AudioType = "ec-3"
	audioNoAudioDescription AudioType = "noAd"

	audioLangPTBR AudioLanguage = "pt-BR"
	audioLangES   AudioLanguage = "es-MX"
	audioLangEN   AudioLanguage = "en"

	captionPTBR CaptionLanguage = "pt-BR"
	captionES   CaptionLanguage = "es-MX"
	captionEN   CaptionLanguage = "en"

	codecHDR10              Codec = "hdr10"
	codecDolbyVision        Codec = "dovi"
	codecHEVC               Codec = "hevc"
	codecH264               Codec = "avc"
	codecAAC                Codec = "aac"
	codecAC3                Codec = "ac-3"
	codecEnhancedAC3        Codec = "ec-3"
	codecNoAudioDescription Codec = "noAd"

	// ProtocolHLS for manifest in hls
	ProtocolHLS Protocol = "hls"
	// ProtocolDASH for manifests in dash
	ProtocolDASH Protocol = "dash"
)

// MediaFilters is a struct that carry all the information passed via url
type MediaFilters struct {
	Videos             []VideoType       `json:",omitempty"`
	VideoSubFilters    Subfilters        `json:",omitempty"`
	Audios             []AudioType       `json:",omitempty"`
	AudioSubFilters    Subfilters        `json:",omitempty"`
	AudioLanguages     []AudioLanguage   `json:",omitempty"`
	CaptionLanguages   []CaptionLanguage `json:",omitempty"`
	CaptionTypes       []CaptionType     `json:",omitempty"`
	FilterStreamTypes  []StreamType      `json:",omitempty"`
	FilterBitrateTypes []StreamType      `json:",omitempty"`
	MaxBitrate         int               `json:",omitempty"`
	MinBitrate         int               `json:",omitempty"`
	Protocol           Protocol          `json:"protocol"`
}

type Subfilters struct {
	MinBitrate int     `json:",omitempty"`
	MaxBitrate int     `json:",omitempty"`
	Codecs     []Codec `json:",omitempty"`
}

var urlParseRegexp = regexp.MustCompile(`(.*?)\((.*)\)`)

// URLParse will generate a MediaFilters struct with
// all the filters that needs to be applied to the
// master manifest. It will also return the master manifest
// url without the filters.
func URLParse(urlpath string) (string, *MediaFilters, error) {
	mf := new(MediaFilters)
	parts := strings.Split(urlpath, "/")
	re := urlParseRegexp
	masterManifestPath := "/"

	if strings.Contains(urlpath, ".m3u8") {
		mf.Protocol = ProtocolHLS
	} else if strings.Contains(urlpath, ".mpd") {
		mf.Protocol = ProtocolDASH
	}

	// set bitrate defaults
	//mf.MinBitrate = 0
	//mf.MaxBitrate = math.MaxInt32
	mf.initializeBitrateRange()

	for _, part := range parts {
		// FindStringSubmatch should return a slice with
		// the full string, the key and filters (3 elements).
		// If it doesn't match, it means that the path is part
		// of the official manifest path so we concatenate to it.
		subparts := re.FindStringSubmatch(part)
		if len(subparts) != 3 {
			masterManifestPath = path.Join(masterManifestPath, part)
			continue
		}

		filters := strings.Split(subparts[2], ",")

		tesRX := regexp.MustCompile(`\),`)
		subfilters := SplitAfter(subparts[2], tesRX)

		switch key := subparts[1]; key {
		case "v":
			for _, sf := range subfilters {
				splitSubfilter := re.FindStringSubmatch(sf)
				var key string
				var param []string
				if len(splitSubfilter) == 0 {
					key = "c"
					param = strings.Split(sf, ",")
					for _, videoType := range filters {
						if videoType == "hdr10" {
							mf.Videos = append(mf.Videos, VideoType("hev1.2"), VideoType("hvc1.2"))
							continue
						}

						mf.Videos = append(mf.Videos, VideoType(videoType))
					}
				} else {
					key = splitSubfilter[1]
					param = strings.Split(splitSubfilter[2], ",")
				}
				mf.normalizeSubfilter(StreamType("video"), key, param)
			}

		case "a":
			for _, sf := range subfilters {
				splitSubfilter := re.FindStringSubmatch(sf) // right now, assuming the nested filters are of valid form, maybe add a check here for that
				var key string
				var param []string
				if len(splitSubfilter) == 0 {
					key = "c"
					param = strings.Split(sf, ",")
					for _, audioType := range filters {
						mf.Audios = append(mf.Audios, AudioType(audioType))
					}
				} else {
					key = splitSubfilter[1]
					param = strings.Split(splitSubfilter[2], ",")
				}
				mf.normalizeSubfilter(StreamType("audio"), key, param)
			}

		case "al":
			for _, audioLanguage := range filters {
				mf.AudioLanguages = append(mf.AudioLanguages, AudioLanguage(audioLanguage))
			}
		case "c":
			for _, captionLanguage := range filters {
				mf.CaptionLanguages = append(mf.CaptionLanguages, CaptionLanguage(captionLanguage))
			}
		case "ct":
			if mf.CaptionTypes == nil {
				mf.CaptionTypes = []CaptionType{}
			}

			for _, captionType := range filters {
				mf.CaptionTypes = append(mf.CaptionTypes, CaptionType(captionType))
			}
		case "fs":
			for _, streamType := range filters {
				mf.FilterStreamTypes = append(mf.FilterStreamTypes, StreamType(streamType))
			}
		case "b":
			if filters[0] != "" {
				mf.MinBitrate, _ = strconv.Atoi(filters[0])
			}

			if filters[1] != "" {
				mf.MaxBitrate, _ = strconv.Atoi(filters[1])
			}
		}
	}

	return masterManifestPath, mf, nil
}

func (f *MediaFilters) initializeBitrateRange() {
	f.MinBitrate = 0
	f.MaxBitrate = math.MaxInt32
	f.AudioSubFilters.MinBitrate = 0
	f.AudioSubFilters.MaxBitrate = math.MaxInt32
	f.VideoSubFilters.MinBitrate = 0
	f.VideoSubFilters.MaxBitrate = math.MaxInt32
}

// SplitAfter splits a string after the matchs of the specified regexp
func SplitAfter(s string, re *regexp.Regexp) []string {
	var splitResults []string
	var position int
	indices := re.FindAllStringIndex(s, -1)
	if indices == nil {
		return append(splitResults, s)
	}
	for _, idx := range indices {
		section := s[position:idx[1]]
		splitResults = append(splitResults, section)
		position = idx[1]
	}
	return append(splitResults, s[position:])
}

// normalizeSubfilter takes a subfilter and sets AudioSubFilters' or VideoSubFilters' values accordingly.
func (f *MediaFilters) normalizeSubfilter(streamType StreamType, key string, values []string) {
	var streamSubfilters *Subfilters
	switch streamType {
	case "audio":
		streamSubfilters = &f.AudioSubFilters
	case "video":
		streamSubfilters = &f.VideoSubFilters
	}

	switch key {
	case "c":
		for _, v := range values {
			if v == "hdr10" {
				streamSubfilters.Codecs = append(streamSubfilters.Codecs, Codec("hev1.2"), Codec("hvc1.2"))
			} else {
				streamSubfilters.Codecs = append(streamSubfilters.Codecs, Codec(v))
			}
		}
	case "b":
		if values[0] != "" {
			streamSubfilters.MinBitrate, _ = strconv.Atoi(values[0])
		}

		if values[1] != "" {
			streamSubfilters.MaxBitrate, _ = strconv.Atoi(values[1])
		}
	}
}
