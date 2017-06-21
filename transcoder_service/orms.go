package main

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

// GetVideoObjects returns a list of Video objects and its count from database
func GetVideoObjects(connection *gorm.DB) (uint, []Video, error) {
	var videos []Video
	var count uint
	var dbError error

	defer connection.Close()

	connection.Find(&videos).Count(&count)
	if connection.Error != nil {
		dbError = connection.Error
	}

	return count, videos, dbError
}

// GetVideoObject returns a Video object from given id from database
func GetVideoObject(videoID int, connection *gorm.DB) (Video, error) {
	var video Video
	var dbError error

	defer connection.Close()

	fmt.Println(video)

	connection.Where(map[string]interface{}{"id": videoID}).First(&video)
	if connection.Error != nil {
		dbError = connection.Error
	}

	if video.ID == 0 {
		dbError = errors.New("no video found")
	}

	return video, dbError
}

// CreateVideoObject pushes Video object to database
func CreateVideoObject(videoSerializer Video, connection *gorm.DB) (Video, error) {
	var dbError error

	defer connection.Close()

	connection.NewRecord(videoSerializer)
	connection.Create(&videoSerializer)
	if connection.Error != nil {
		dbError = connection.Error
	}

	return videoSerializer, dbError
}

// CreateVideoRenderingObject pushes VideoRendering object to database
func CreateVideoRenderingObject(videoRendering VideoRendering, connection *gorm.DB) (VideoRendering, error) {
	var dbError error

	defer connection.Close()

	connection.NewRecord(videoRendering)
	connection.Create(&videoRendering)
	if connection.Error != nil {
		dbError = connection.Error
	}

	return videoRendering, dbError
}
