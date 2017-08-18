package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// GetDatabaseConnection returns an instance
// as a database connection
func GetDatabaseConnection(user string, password string, host string, db string) *gorm.DB {
	connection, err := gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, db))

	if err != nil {
		panic(err)
	}

	return connection
}

// CreateSchemas creates a set of database tables
// from Go struct classes
func CreateSchemas(user string, password string, host string, db string) {
	connection := GetDatabaseConnection(user, password, host, db)

	defer connection.Close()

	connection.AutoMigrate(&Video{}, &VideoRendering{})

	connection.Model(&VideoRendering{}).AddForeignKey("video_id", "videos(id)", "CASCADE", "CASCADE")

	connection.Model(&VideoRendering{}).AddIndex("idx_video_id", "video_id")

}
