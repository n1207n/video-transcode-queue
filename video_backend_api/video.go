package main

import (
	"fmt"
	"time"
)

// Video represents an uploaded video instance
// Relation:
// - has many VideoRendering
type Video struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	Title          string `gorm:"not null" json:"title" binding:"required"`
	IsReadyToServe bool   `sql:"DEFAULT:false" json:"is_ready_to_serve"`
	StreamFilePath string `json:"stream_file_path"`

	Renderings []VideoRendering `gorm:"ForeignKey:VideoID"`
}

func (v Video) String() string {
	return fmt.Sprintf("Video: %d - %s", v.ID, v.Title)
}

// VideoRendering represents each rendering variant from original
type VideoRendering struct {
	ID             uint      `gorm:"primary_key" json:"id"`
	CreatedAt      time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"not null" json:"updated_at"`
	RenderingTitle string    `gorm:"not null" json:"rendering_title" binding:"required"`

	FilePath string `gorm:"not null" json:"file_path"`
	URL      string `gorm:"not null" json:"url"`
	Width    uint   `gorm:"not null" json:"width"`
	Height   uint   `gorm:"not null" json:"height"`

	VideoID uint `json:"video_id"`
}

func (vr VideoRendering) String() string {
	return fmt.Sprintf("VideoRendering: %d - %s", vr.ID, vr.RenderingTitle)
}
