package filters

import (
	"math"
	"strings"

	"github.com/cbsinteractive/bakery/pkg/parsers"
)

// Filter is an interface for HLS and DASH filters
type Filter interface {
	FilterManifest(filters *parsers.MediaFilters) (string, error)
}

// ContentType represents the content in the stream
type ContentType string

const (
	captionContentType ContentType = "text"
	audioContentType   ContentType = "audio"
	videoContentType   ContentType = "video"
)

// CodecFilterID is the formatted codec represented in a given playlist
type CodecFilterID string

const (
	hevcCodec  CodecFilterID = "hvc"
	avcCodec   CodecFilterID = "avc"
	dolbyCodec CodecFilterID = "dvh"
	aacCodec   CodecFilterID = "mp4a"
	ec3Codec   CodecFilterID = "ec-3"
	ac3Codec   CodecFilterID = "ac-3"
	stppCodec  CodecFilterID = "stpp"
	wvttCodec  CodecFilterID = "wvtt"
)

// ValidCodecs returns a map of all formatted values for a given codec filter
func ValidCodecs(codec string, filter CodecFilterID) bool {
	return strings.Contains(codec, string(filter))
}

// ValidBitrateRange returns true if the specified min and max bitrates create a valid range
func ValidBitrateRange(minBitrate int, maxBitrate int) bool {
	return (minBitrate >= 0 && maxBitrate <= math.MaxInt32) &&
		(minBitrate < maxBitrate) &&
		!(minBitrate == 0 && maxBitrate == math.MaxInt32)
}

// DefinesBitrateFilter returns true if a bitrate filter should be applied. This means that
// at least one of the overall, audio, and video bitrate filters are valid and not the default range
func DefinesBitrateFilter(f *parsers.MediaFilters) bool {
	overall := ValidBitrateRange(f.MinBitrate, f.MaxBitrate)
	if overall {
		f.AudioSubFilters.MinBitrate = max(f.AudioSubFilters.MinBitrate, f.MinBitrate)
		f.AudioSubFilters.MaxBitrate = min(f.AudioSubFilters.MaxBitrate, f.MaxBitrate)
		f.VideoSubFilters.MinBitrate = max(f.VideoSubFilters.MinBitrate, f.MinBitrate)
		f.VideoSubFilters.MaxBitrate = min(f.VideoSubFilters.MaxBitrate, f.MaxBitrate)
		return true
	} else {
		audio := ValidBitrateRange(f.AudioSubFilters.MinBitrate, f.AudioSubFilters.MaxBitrate)
		video := ValidBitrateRange(f.VideoSubFilters.MinBitrate, f.VideoSubFilters.MaxBitrate)
		return audio || video
	}
}

// max returns the larger of int a and int b
func max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

// min returns the smaller of int a and int b
func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
