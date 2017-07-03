package main

import (
	"os"

	"go.uber.org/zap"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
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
	sugaredLogger.Info("Starting video streaming API server")

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	// TODO: Use uploadFolderPath later for better security
	router.Use(static.Serve("/contents", static.LocalFile("/", false)))

	// By default it serves on :8080
	router.Run(":8880")
}
