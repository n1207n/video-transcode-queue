package main

import "github.com/jinzhu/gorm"

// GetVideoObjects returns a list of Video objects and its count from database
func GetVideoObjects(connection *gorm.DB) (uint, []*Video, error) {
	var videos []*Video
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
func GetVideoObject(videoID int, connection *gorm.DB) (*Video, error) {
	var video *Video
	var dbError error

	defer connection.Close()

	connection.First(&video, videoID)
	if connection.Error != nil {
		dbError = connection.Error
	}

	return video, dbError
}

// CreateVideoObject pushes Video object to database
func CreateVideoObject(videoSerializer *Video, connection *gorm.DB) (*Video, error) {
	var dbError error

	defer connection.Close()

	connection.NewRecord(videoSerializer)
	connection.Create(&videoSerializer)
	if connection.Error != nil {
		dbError = connection.Error
	}

	return videoSerializer, dbError
}
