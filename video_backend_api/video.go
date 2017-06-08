package main

// VideoRendering represents each rendering variant from original
type VideoRendering struct {
	tableName struct{} `sql:"video_renderings, alias:video_rendering"`

	ID             string `json:"id"`
	RenderingTitle string `sql:",notnull",json:"rendering_title"`

	FilePath string `sql:",notnull",json:"file_path"`
	URL      string `sql:",notnull",json:"url"`
	Width    uint   `sql:",notnull",json:"width"`
	Height   uint   `sql:",notnull",json:"height"`

	VideoID int
}

// Video represents an uploaded video instance
// Relation:
// - has many VideoRendering
type Video struct {
	tableName struct{} `sql:"videos, alias:video"`

	ID             string `json:"id"`
	Title          string `sql:",notnull",json:"title"`
	IsReadyToServe bool   `sql:",notnull",json:"is_ready_to_serve"`

	Renders []*VideoRendering
}

// VideoCreate represents a video instance to upload
type VideoCreate struct {
	tableName struct{} `sql:"videos, alias:video"`
	Title     string   `sql:",notnull",json:"title"`
}
