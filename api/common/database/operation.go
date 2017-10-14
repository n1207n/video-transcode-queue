package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/n1207n/video-transcode-platform/api/common/entity"
)

// GetConnection returns an instance
// as a database connection
func GetConnection(user string, password string, host string, db string) *gorm.DB {
	connection, err := gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, db))

	if err != nil {
		panic(err)
	}

	return connection
}

// CreateSchemas creates a set of database tables
// from Go struct classes
func CreateSchemas(user string, password string, host string, db string) {
	connection := GetConnection(user, password, host, db)

	defer connection.Close()

	connection.AutoMigrate(&entity.Video{}, &entity.VideoRendering{})

	connection.Model(&entity.VideoRendering{}).AddForeignKey("video_id", "videos(id)", "CASCADE", "CASCADE")

	connection.Model(&entity.VideoRendering{}).AddIndex("idx_video_id", "video_id")

}
