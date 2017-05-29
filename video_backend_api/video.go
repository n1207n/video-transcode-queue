package main

import "github.com/go-pg/pg"

// VideoRendering represents each rendering variant from original
type VideoRendering struct {
	tableName struct{} `sql:"video_renderings, alias:video_rendering"`

	ID             string `json:"id"`
	RenderingTitle string `json:"rendering_title"`

	FilePath string `json:"file_path"`
	URL      string `json:"url"`
	Width    uint   `json:"width"`
	Height   uint   `json:"height"`
}

// Video represents an uploaded video instance
// Relation:
// - has many VideoRendering
type Video struct {
	tableName struct{} `sql:"videos, alias:video"`

	ID             string `json:"id"`
	Title          string `json:"title"`
	IsReadyToServe bool   `json:"is_ready_to_serve"`

	Renders []*VideoRendering
}

// GetVideoObjects returns a list of Video objects and its count from database
func GetVideoObjects() (int, []*Video, error) {
	var videos []*Video
	var db *pg.DB = getDatabaseConnection()
	var dbError error

	defer func() {
		db.Close()
	}()

	count, err := db.Model(&videos).Order("id DESC").SelectAndCount()
	if err != nil {
		dbError = err
	}

	return count, videos, dbError
}

// GetVideoObject returns a Video object from given id from database
func GetVideoObject(videoID int) (*Video, error) {
	var video *Video
	var db *pg.DB = getDatabaseConnection()
	var dbError error

	defer func() {
		db.Close()
	}()

	err := db.Model(&video).
		Where("id = ?", videoID).
		Select()
	if err != nil {
		dbError = err
	}

	return video, dbError
}
