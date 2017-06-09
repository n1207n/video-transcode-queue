package main

// Video represents an uploaded video instance
// Relation:
// - has many VideoRendering
type Video struct {
	tableName struct{} `sql:"videos, alias:video"`

	ID             uint   `json:"id"`
	Title          string `sql:",notnull",json:"title"`
	IsReadyToServe bool   `sql:",notnull",json:"is_ready_to_serve"`

	Renders []*VideoRendering
}

// VideoRendering represents each rendering variant from original
type VideoRendering struct {
	tableName struct{} `sql:"video_renderings, alias:video_rendering"`

	ID             uint   `json:"id"`
	RenderingTitle string `sql:",notnull",json:"rendering_title"`

	FilePath string `sql:",notnull",json:"file_path"`
	URL      string `sql:",notnull",json:"url"`
	Width    uint   `sql:",notnull",json:"width"`
	Height   uint   `sql:",notnull",json:"height"`

	VideoID int
	Video   *Video
}

// VideoCreateAPISerializer represents a video instance to create
type VideoCreateAPISerializer struct {
	tableName struct{} `sql:"videos, alias:video"`

	Title string `json:"title"`
}
