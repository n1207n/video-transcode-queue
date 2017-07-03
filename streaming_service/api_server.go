package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/adjust/rmq"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

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
	sugaredLogger.Info("Starting video streaming API server")

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.Use(static.Serve("/contents", static.LocalFile(uploadFolderPath, false)))

	// By default it serves on :8080
	router.Run(":8880")
}
