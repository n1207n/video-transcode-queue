package main

import (
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/golang/glog"
	"gopkg.in/gin-gonic/gin.v1"
)

var (
	pgDb, pgUser, pgPassword, pgHost string
	queueTopic                       string
)

func main() {
	loadEnvironmentVariables()

	err := createSchemas()
	if err != nil {
		glog.Infoln("Failed to sync up database tables. Check the connection.")
		panic(err)
	}

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

	queueTopic = os.Getenv("QUEUE_TOPIC")
	if len(queueTopic) == 0 {
		panic("No QUEUE_TOPIC environment variable")
	}
}

// getDatabaseConnection returns an instance
// as a database connection
func getDatabaseConnection() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     pgUser,
		Password: pgPassword,
		Database: pgDb,
		Addr:     pgHost,
	})

	return db
}

// createSchemas creates a set of database tables
// from Go struct classes
func createSchemas() error {
	var db *pg.DB = getDatabaseConnection()

	defer func() {
		db.Close()
	}()

	for _, model := range []interface{}{&VideoRendering{}, &Video{}} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func startAPIServer() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	v1 := router.Group("/v1")
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

}

func getVideoDetail(c *gin.Context) {

}

func createVideo(c *gin.Context) {

}

func uploadVideoFile(c *gin.Context) {

}
