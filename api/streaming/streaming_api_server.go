package main

import (
	"os"

	"go.uber.org/zap"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	pgDb, pgUser, pgPassword, pgHost string
	uploadFolderPath                 string
	logger                           *zap.SugaredLogger
)

func main() {
	loadEnvironmentVariables()
	startStreamingAPIServer()
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

func startStreamingAPIServer() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	logger = log.Sugar()
	logger.Info("Starting streaming API server")

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Range"}

	router.Use(cors.New(corsConfig))

	// TODO: Use uploadFolderPath later for better security
	router.Use(static.Serve("/contents", static.LocalFile("/", false)))

	// By default it serves on :8080
	router.Run(":8880")
}
