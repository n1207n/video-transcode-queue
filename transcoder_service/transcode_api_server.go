package main

import (
	"net/http"

	"github.com/golang/glog"
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

		performTranscoding(request.Path)
		c.JSON(http.StatusAccepted, gin.H{"video_id": request.VideoID})
	}
}

// TranscodeRequest represents a JSON POST data for video-transcode API
type TranscodeRequest struct {
	Path    string `json:"path"`
	VideoID string `json:"video_id"`
}

func performTranscoding(filePath string) {

}
