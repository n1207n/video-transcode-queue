package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang/glog"
	"gopkg.in/gin-gonic/gin.v1"
)

var (
	pgDb, pgUser, pgPassword, pgHost string
	uploadFolderPath                 string
)

func main() {
	loadEnvironmentVariables()
	glog.Infoln("Starting transcoder API server")
	startAPIServer()
}

// loadEnvironmentVariables loads PostgreSQL
// information from dotenv
func loadEnvironmentVariables() {
	pgDb = os.Getenv("PGDB")
	if len(pgDb) == 0 {
		panic("No PGDB environment variable")
	}

	pgUser = os.Getenv("PGUSER")
	if len(pgUser) == 0 {
		panic("No PGUSER environment variable")
	}

	pgPassword = os.Getenv("PGPASSWORD")
	if len(pgPassword) == 0 {
		panic("No PGPASSWORD environment variable")
	}

	pgHost = os.Getenv("PGHOST")
	if len(pgHost) == 0 {
		panic("No PGHOST environment variable")
	}

	uploadFolderPath = os.Getenv("UPLOAD_FOLDER_PATH")
	if len(uploadFolderPath) == 0 {
		panic("No UPLOAD_FOLDER_PATH environment variable")
	}
}

func startAPIServer() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	v1 := router.Group("/api/v1")
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
