package main

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/transcode"
	"github.com/nareix/joy4/cgo/ffmpeg"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	glog.Infoln("Starting transcoder API server")
	startAPIServer()
}

func startAPIServer() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.POST("/video-transcode", transcodeVideo)
	}

	// By default it serves on :8080
	router.Run()
}

func transcodeVideo(c *gin.Context) {
	var request TranscodeRequest

	if c.BindJSON(&request) == nil {
		// TODO: Check if video id exists in PostgreSQL

		transcodedFilePaths, err := performTranscoding(request.Path)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to open a video file: %s for Video ID %s", request.Path, request.VideoID)})
		} else {
			c.JSON(http.StatusAccepted, gin.H{"video_id": request.VideoID})
		}
	}
}

// TranscodeRequest represents a JSON POST data for video-transcode API
type TranscodeRequest struct {
	Path    string `json:"path"`
	VideoID string `json:"video_id"`
}

func performTranscoding(filePath string) (transcodedFilePaths []string, transcodeError error) {
	infile, err := avutil.Open(filePath)

	if err != nil {
		transcodeError = err
	}

	encoded_audio_info := func(stream av.AudioCodecData, i int) (need bool, dec av.AudioDecoder, enc av.AudioEncoder, err error) {
		need = true
		dec, _ = ffmpeg.NewAudioDecoder(stream)
		enc, _ = ffmpeg.NewAudioEncoderByName("libfdk_aac")
		enc.SetSampleRate(stream.SampleRate())
		enc.SetChannelLayout(av.CH_STEREO)
		enc.SetBitrate(12000)
		enc.SetOption("profile", "HE-AACv2")
		return
	}

	transcode_demuxer := &transcode.Demuxer{
		Options: transcode.Options{
			FindAudioDecoderEncoder: encoded_audio_info,
		},
		Demuxer: infile,
	}

	outfile, err := avutil.Create("output.ts")

	if err != nil {
		transcodeError = err
	} else {
		avutil.CopyFile(outfile, trans)

		outfile.Close()
		infile.Close()
		transcode_demuxer.Close()
	}

	return
}
