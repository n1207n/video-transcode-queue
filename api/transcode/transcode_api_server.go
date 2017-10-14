package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

var (
	pgDb, pgUser, pgPassword, pgHost string
	uploadFolderPath                 string
	logger                           *zap.SugaredLogger
)

func main() {
	loadEnvironmentVariables()
	startTranscodeAPIServer()
}

// loadEnvironmentVariables loads PostgreSQL
// information from dotenv
func loadEnvironmentVariables() {
	pgDb = os.Getenv("PGDB")
	if len(pgDb) == 0 {
		panic("No pgDB environment variable")
	}

	pgUser = os.Getenv("PGUSER")
	if len(pgUser) == 0 {
		panic("No pgUSER environment variable")
	}

	pgPassword = os.Getenv("PGPASSWORD")
	if len(pgPassword) == 0 {
		panic("No pgPASSWORD environment variable")
	}

	pgHost = os.Getenv("PGHOST")
	if len(pgHost) == 0 {
		panic("No pgHOST environment variable")
	}

	uploadFolderPath = os.Getenv("UPLOAD_FOLDER_PATH")
	if len(uploadFolderPath) == 0 {
		panic("No UPLOAD_FOLDER_PATH environment variable")
	}
}

func startTranscodeAPIServer() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	logger = log.Sugar()
	logger.Info("Starting transcode API server")

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

func performTranscoding(request TranscodeRequest, c *gin.Context) (transcodedFilePaths []string, transcodeError error) {
	splitStringPaths := strings.Split(request.Path, "/")
	fileFolderPath := strings.Join(splitStringPaths[:len(splitStringPaths)-1], "/")
	filename := splitStringPaths[len(splitStringPaths)-1]

	// Strip the file extension and convert any reverse subsequent dots to underscore
	splitFilenameCharacters := strings.Split(filename, ".")
	videoName := strings.Join(splitFilenameCharacters[:len(splitFilenameCharacters)-1], "_")

	videoID, _ := strconv.Atoi(request.VideoID)

	var waitGroup sync.WaitGroup

	_, height, err := GetVideoDimensionInfo(filename, fileFolderPath, logger)
	if err != nil {
		logger.Errorf("Error from getting video dimension info: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, gin.H{"video_id": request.VideoID, "status": "Failed to get video metadata. Corrupted file?"})

		// TODO: Delete the video file
	} else {
		var targets []int
		dbConnectionInfo := map[string]string{
			"pgDb":       pgDb,
			"pgUser":     pgUser,
			"pgPassword": pgPassword,
			"pgHost":     pgHost,
		}

		if height >= 720 {
			targets = append(targets, 720)
		}

		if height >= 540 {
			targets = append(targets, 540)
		}

		if height >= 360 {
			targets = append(targets, 360)
		}

		if height < 360 {
			targets = append(targets, 360)
		}

		waitGroup.Add(len(targets))

		for _, target := range targets {
			switch target {
			case 720:
				go TranscodeToHD720P(videoName, videoID, filename, fileFolderPath, dbConnectionInfo, &waitGroup, logger)
			case 540:
				go TranscodeToSD540P(videoName, videoID, filename, fileFolderPath, dbConnectionInfo, &waitGroup, logger)
			case 360:
				go TranscodeToSD360P(videoName, videoID, filename, fileFolderPath, dbConnectionInfo, &waitGroup, logger)
			default:
				go TranscodeToSD360P(videoName, videoID, filename, fileFolderPath, dbConnectionInfo, &waitGroup, logger)
			}
		}

		waitGroup.Wait()

		c.JSON(http.StatusOK, gin.H{"video_id": request.VideoID, "status": "In progress"})

		logger.Infof("Constructing MPD for %s", videoName)

		ConstructMPD(videoName, videoID, filename, fileFolderPath, targets, dbConnectionInfo, logger)
	}

	return
}
