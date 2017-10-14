package entity

import (
	"time"
)

// Task represents a job task to transcode a video file
type Task struct {
	ID        string
	FilePath  string
	Timestamp time.Time
}
