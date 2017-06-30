package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"go.uber.org/zap"

	"gopkg.in/gin-gonic/gin.v1"
)

var (
	pgDb, pgUser, pgPassword, pgHost string
	uploadFolderPath                 string
	sugaredLogger                    *zap.SugaredLogger
)

func main() {
	loadEnvironmentVariables()
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
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugaredLogger = logger.Sugar()
	sugaredLogger.Info("Starting transcoder API server\n")

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/video-transcode", transcodeVideo)
	}

	// By default it serves on :8080
	router.Run(":8800")
}

func transcodeVideo(c *gin.Context) {
	var request TranscodeRequest

	if c.BindJSON(&request) == nil {
		// TODO: Check if video id exists in PostgreSQL

		performTranscoding(request, c)
	}
}

// TranscodeRequest represents a JSON POST data for video-transcode API
type TranscodeRequest struct {
	Path    string `json:"path" binding:"required"`
	VideoID string `json:"video_id" binding:"required"`
}

func performTranscoding(transcodeRequest TranscodeRequest, c *gin.Context) (transcodedFilePaths []string, transcodeError error) {
	splitStringPaths := strings.Split(transcodeRequest.Path, "/")
	fileFolderPath := strings.Join(splitStringPaths[:len(splitStringPaths)-1], "/")
	filename := splitStringPaths[len(splitStringPaths)-1]

	// Strip the file extension and convert any reverse subsequent dots to underscore
	splitFilenameCharacters := strings.Split(filename, ".")
	videoName := strings.Join(splitFilenameCharacters[:len(splitFilenameCharacters)-1], "_")

	videoID, _ := strconv.Atoi(transcodeRequest.VideoID)

	var waitGroup sync.WaitGroup

	_, height, err := GetVideoDimensionInfo(filename, fileFolderPath)
	if err != nil {
		sugaredLogger.Errorf("Error from getting video dimension info: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, gin.H{"video_id": transcodeRequest.VideoID, "status": "Failed to get video metadata. Corrupted file?"})

		// TODO: Delete the video file
	} else {
		var transcodeTargets []int
		dbConnectionInfo := map[string]string{
			pgDb:       pgDb,
			pgUser:     pgUser,
			pgPassword: pgPassword,
			pgHost:     pgHost,
		}

		if height >= 720 {
			transcodeTargets = append(transcodeTargets, 720)
		}

		if height >= 540 {
			transcodeTargets = append(transcodeTargets, 540)
		}

		if height >= 360 {
			transcodeTargets = append(transcodeTargets, 360)
		}

		if height < 360 {
			transcodeTargets = append(transcodeTargets, 360)
		}

		waitGroup.Add(len(transcodeTargets))

		for _, target := range transcodeTargets {
			switch target {
			case 720:
				go TranscodeToHD720P(videoName, videoID, filename, fileFolderPath, dbConnectionInfo, &waitGroup)
			case 540:
				go TranscodeToSD540P(videoName, videoID, filename, fileFolderPath, dbConnectionInfo, &waitGroup)
			case 360:
				go TranscodeToSD360P(videoName, videoID, filename, fileFolderPath, dbConnectionInfo, &waitGroup)
			default:
				go TranscodeToSD360P(videoName, videoID, filename, fileFolderPath, dbConnectionInfo, &waitGroup)
			}
		}

		waitGroup.Wait()

		c.JSON(http.StatusOK, gin.H{"video_id": transcodeRequest.VideoID, "status": "In progress"})

		sugaredLogger.Infof("Constructing MPD for %s", videoName)

		ConstructMPD(videoName, videoID, filename, fileFolderPath, transcodeTargets, dbConnectionInfo)
	}

	return
}
