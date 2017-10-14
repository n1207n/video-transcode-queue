package database

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/n1207n/video-transcode-queue/api/common/entity"
)

// GetVideoObjects returns a list of Video objects and its count from database
func GetVideoObjects(connection *gorm.DB) (uint, []entity.Video, error) {
	var videos []entity.Video
	var count uint
	var dbError error

	defer connection.Close()

	connection.Preload("Renderings").Find(&videos)
	if connection.Error != nil {
		dbError = connection.Error
	}

	return count, videos, dbError
}

// GetVideoRenderingObjects returns a list of VideoRendering objects linked to Video ID and its count from database
func GetVideoRenderingObjects(video entity.Video, connection *gorm.DB) (uint, []entity.VideoRendering, error) {
	var renderings []entity.VideoRendering
	var count uint
	var dbError error

	defer connection.Close()

	connection.Model(&video).Related(&renderings).Count(&count)
	if connection.Error != nil {
		dbError = connection.Error
	}

	return count, renderings, dbError
}

// GetVideoObject returns a Video object from given id from database
func GetVideoObject(videoID int, connection *gorm.DB) (entity.Video, error) {
	var video entity.Video
	var dbError error

	defer connection.Close()

	connection.Where(map[string]interface{}{"id": videoID}).Preload("Renderings").First(&video)
	if connection.Error != nil {
		dbError = connection.Error
	}

	if video.ID == 0 {
		dbError = errors.New("no video found")
	}

	return video, dbError
}

// CreateVideoObject pushes Video object to database
func CreateVideoObject(videoSerializer entity.Video, connection *gorm.DB) (entity.Video, error) {
	var dbError error

	defer connection.Close()

	connection.NewRecord(videoSerializer)
	connection.Create(&videoSerializer)
	if connection.Error != nil {
		dbError = connection.Error
	}

	return videoSerializer, dbError
}

// UpdateVideoObject updates Video object to database
func UpdateVideoObject(updatedVideo entity.Video, connection *gorm.DB) (entity.Video, error) {
	var dbError error

	defer connection.Close()

	connection.Save(&updatedVideo)
	if connection.Error != nil {
		dbError = connection.Error
	}

	return updatedVideo, dbError
}

// DeleteVideoObject deletes Video object in database
func DeleteVideoObject(video entity.Video, connection *gorm.DB) (entity.Video, error) {
	var dbError error

	defer connection.Close()

	connection.Delete(&video)
	if connection.Error != nil {
		dbError = connection.Error
	}

	return video, dbError
}

// CreateVideoRenderingObject pushes VideoRendering object to database
func CreateVideoRenderingObject(videoRendering entity.VideoRendering, connection *gorm.DB) (entity.VideoRendering, error) {
	var dbError error

	defer connection.Close()

	connection.NewRecord(videoRendering)
	connection.Create(&videoRendering)
	if connection.Error != nil {
		dbError = connection.Error
	}

	return videoRendering, dbError
}
