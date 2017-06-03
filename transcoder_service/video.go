package main

// VideoRendering represents each rendering variant from original
type VideoRendering struct {
	tableName struct{} `sql:"video_renderings, alias:video_rendering"`

	ID             string `json:"id"`
	RenderingTitle string `json:"rendering_title"`

	FilePath string `json:"file_path"`
	URL      string `json:"url"`
	Width    uint   `json:"width"`
	Height   uint   `json:"height"`

	VideoID int
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
