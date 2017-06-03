package main

import "github.com/go-pg/pg"

// GetVideoObjects returns a list of Video objects and its count from database
func GetVideoObjects(connection *pg.DB) (int, []*Video, error) {
	var videos []*Video
	var dbError error

	defer func() {
		connection.Close()
	}()

	count, err := connection.Model(&videos).Order("id DESC").SelectAndCount()
	if err != nil {
		dbError = err
	}

	return count, videos, dbError
}

// GetVideoObject returns a Video object from given id from database
func GetVideoObject(videoID int, connection *pg.DB) (*Video, error) {
	var video *Video
	var dbError error

	defer func() {
		connection.Close()
	}()

	err := connection.Model(&video).
		Where("id = ?", videoID).
		Select()
	if err != nil {
		dbError = err
	}

	return video, dbError
}
