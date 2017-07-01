package main

import "time"

// FFProbeStreamData represents JSON format for each stream
type FFProbeStreamData struct {
	Index              int               `json:"index"`
	CodecName          string            `json:"codec_name"`
	CodecLongName      string            `json:"codec_long_name"`
	Profile            int               `json:"profile,string"`
	CodecType          string            `json:"codec_type"`
	CodecTimeBase      string            `json:"codec_time_base"`
	CodecTagString     string            `json:"codec_tag_string"`
	CodecTag           string            `json:"codec_tag"`
	Width              *int              `json:"width,omitempty"`
	Height             *int              `json:"height,omitempty"`
	CodedWidth         *int              `json:"coded_width,omitempty"`
	CodedHeight        *int              `json:"coded_height,omitempty"`
	HasBFrames         *int              `json:"has_b_frames,omitempty"`
	SampleAspectRatio  *string           `json:"sample_aspect_ratio,omitempty"`
	DisplayAspectRatio *string           `json:"display_aspect_ratio,omitempty"`
	PixFmt             *string           `json:"pix_fmt,omitempty"`
	Level              *int              `json:"level,omitempty"`
	ColorRange         *string           `json:"color_range,omitempty"`
	ColorSpace         *string           `json:"color_space,omitempty"`
	ColorTransfer      *string           `json:"color_transfer,omitempty"`
	ColorPrimaries     *string           `json:"color_primaries,omitempty"`
	ChromaLocation     *string           `json:"chroma_location,omitempty"`
	Refs               *int              `json:"refs,omitempty"`
	IsAVC              *bool             `json:"is_avc,string,omitempty"`
	NalLengthSize      *int              `json:"nal_length_size,string,omitempty"`
	RFrameRate         string            `json:"r_frame_rate"`
	AVGFrameRate       string            `json:"avg_frame_rate"`
	TimeBase           string            `json:"time_base"`
	StartPTS           int               `json:"start_pts"`
	StartTime          float64           `json:"start_time,string"`
	DurationTS         int               `json:"duration_ts"`
	Duration           float64           `json:"duration,string"`
	BitRate            int               `json:"bit_rate,string"`
	BitsPerRawSample   *int              `json:"bits_per_raw_sample,string,omitempty"`
	NBFrames           int               `json:"nb_frames,string"`
	SampleFMT          *string           `json:"sample_fmt,omitempty"`
	SampleRate         *int              `json:"sample_rate,string,omitempty"`
	Channels           *int              `json:"channels,omitempty"`
	ChannelLayout      *string           `json:"channel_layout,omitempty"`
	BitsPerSample      *int              `json:"bits_per_sample,omitempty"`
	Disposition        map[string]int    `json:"disposition"`
	Tags               map[string]string `json:"tags"`
}

// StartTimeDuration represents
// FFProbeStreamData's StartTime field as Duration object
func (f FFProbeStreamData) StartTimeDuration() time.Duration {
	return time.Duration(f.StartTime * float64(time.Second))
}

// DurationAsObject represents
// FFProbeStreamData's Duration field as Duration object
func (f FFProbeStreamData) DurationAsObject() time.Duration {
	return time.Duration(f.Duration * float64(time.Second))
}

// ProbeData represents ffprobe info as JSON struct
type ProbeData struct {
	Stream []FFProbeStreamData `json:"streams,omitempty"`
}
