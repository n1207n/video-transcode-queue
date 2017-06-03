package main

import (
	"net/http"
	"strings"

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

		c.JSON(http.StatusAccepted, gin.H{"video_id": request.VideoID, "status": "In progress"})
	}
}

// TranscodeRequest represents a JSON POST data for video-transcode API
type TranscodeRequest struct {
	Path    string `json:"path"`
	VideoID string `json:"video_id"`
}

func performTranscoding(filePath string) (transcodedFilePaths []string, transcodeError error) {
	splitStringPaths := strings.Split(filePath, "/")
	fileFolderPath := strings.Join(splitStringPaths[:len(splitStringPaths)-1], "/")
	filename := splitStringPaths[len(splitStringPaths)-1]

	// Strip the file extension and convert any reverse subsequent dots to underscore
	splitFilenameCharacters := strings.Split(filename, ".")
	videoName := strings.Join(splitFilenameCharacters[:len(splitFilenameCharacters)-1], "_")

	go TranscodeToStandard(videoName, filename, fileFolderPath)
	go TranscodeToMobile(videoName, filename, fileFolderPath)
	go TranscodeToHighSD(videoName, filename, fileFolderPath)

	return
}
