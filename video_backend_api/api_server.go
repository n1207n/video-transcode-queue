package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/adjust/rmq"

	"gopkg.in/gin-gonic/gin.v1"
	redis "gopkg.in/redis.v3"
)

var (
	pgDb, pgUser, pgPassword, pgHost               string
	uploadFolderPath                               string
	redisURL, redisPort, redisPassword, redisTopic string
	redisProtocol                                  = "tcp"
	redisNetworkTag                                = "transcode_task_consume"
	sugaredLogger                                  *zap.SugaredLogger
)

func main() {
	loadEnvironmentVariables()
	CreateSchemas(pgUser, pgPassword, pgHost, pgDb)
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

	redisURL = os.Getenv("REDIS_URL")
	if len(redisURL) == 0 {
		panic("No REDIS_URL environment variable")
	}

	redisPort = os.Getenv("REDIS_PORT")
	if len(redisPort) == 0 {
		panic("No REDIS_PORT environment variable")
	}

	redisPassword = os.Getenv("REDIS_PASSWORD")
	if len(redisPassword) == 0 {
		panic("No REDIS_PASSWORD environment variable")
	}

	redisTopic = os.Getenv("REDIS_TOPIC")
	if len(redisTopic) == 0 {
		panic("No REDIS_TOPIC environment variable")
	}
}

// openTaskQueue connects to redis and return a Queue interface
func openTaskQueue() rmq.Queue {
	redisClient := redis.NewClient(&redis.Options{
		Network:  redisProtocol,
		Addr:     fmt.Sprintf("%s:%s", redisURL, redisPort),
		DB:       int64(1),
		Password: redisPassword,
	})

	connection := rmq.OpenConnectionWithRedisClient(redisNetworkTag, redisClient)

	sugaredLogger.Infof("Connected to Redis task queue: %s\n", connection.Name)

	return connection.OpenQueue(redisTopic)
}

func startAPIServer() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugaredLogger = logger.Sugar()
	sugaredLogger.Info("Starting video API server")

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/videos", getVideoList)
		v1.GET("/videos/:id", getVideoDetail)
		v1.POST("/videos", createVideo)
		v1.POST("/video-upload", uploadVideoFile)
	}

	// By default it serves on :8080
	router.Run()
}

func getVideoList(c *gin.Context) {
	connection := GetDatabaseConnection(pgUser, pgPassword, pgHost, pgDb)
	count, videos, err := GetVideoObjects(connection)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"count":   count,
			"results": videos,
		})
	}
}

func getVideoDetail(c *gin.Context) {
	videoID, err := strconv.Atoi(c.Param("id"))

	connection := GetDatabaseConnection(pgUser, pgPassword, pgHost, pgDb)
	video, err := GetVideoObject(videoID, connection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": video,
		})
	}
}

func createVideo(c *gin.Context) {
	var videoSerializer Video

	if err := c.BindJSON(&videoSerializer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	connection := GetDatabaseConnection(pgUser, pgPassword, pgHost, pgDb)
	videoSerializer, err := CreateVideoObject(videoSerializer, connection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"title":   videoSerializer.Title,
		"message": "Object created. Please upload the file for this Video.",
	})
}

func uploadVideoFile(c *gin.Context) {
	c.Request.ParseMultipartForm(64 << 25)

	videoID := c.PostForm("video_id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "video_id is required",
		})

		return
	}

	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Please upload file with 'upload' form field key.",
		})

		return
	}

	filename := header.Filename

	videoFolderPath := uploadFolderPath + videoID
	os.MkdirAll(videoFolderPath, os.ModePerm)

	videoFullPath := fmt.Sprintf("%s/%s", videoFolderPath, filename)

	outFile, err := os.Create(videoFullPath)
	if err != nil {
		sugaredLogger.Fatal("Failed to write filesystem:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "File upload is having issues right now. Please try later.",
		})

		return
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		sugaredLogger.Fatal("Failed to copy video file:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "File upload is having issues right now. Please try later.",
		})

		return
	}

	taskQueue := openTaskQueue()
	task := Task{ID: videoID, Timestamp: time.Now(), FilePath: videoFullPath}

	queueDataBytes, err := json.Marshal(task)
	taskQueue.PublishBytes(queueDataBytes)
	sugaredLogger.Info("Queue task created...:", task)

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Video file uploaded. Transcoding now: %s", videoID),
	})
}
