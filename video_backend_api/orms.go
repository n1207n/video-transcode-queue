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
func GetVideoObject(videoID string, connection *pg.DB) (*Video, error) {
	var video *Video
	var dbError error

	defer func() {
		connection.Close()
	}()

	err := connection.Model(&Video{}).
		Where("id = ?", videoID).
		Select(video)
	if err != nil {
		dbError = err
	}

	return video, dbError
}

// CreateVideoObject pushes Video object to database
func CreateVideoObject(video *Video, connection *pg.DB) (*Video, error) {
	var dbError error

	defer func() {
		connection.Close()
	}()

	err := connection.Insert(&video)
	if err != nil {
		dbError = err
	}

	return video, dbError
}
