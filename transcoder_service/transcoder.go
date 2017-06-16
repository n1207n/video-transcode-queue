package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/golang/glog"
)

// ExecuteFfmpegCLI executes constructed command string by os.exec.Command
func ExecuteFfmpegCLI(commandString string) error {
	commandArguments := strings.Fields(commandString)
	head, commandArguments := commandArguments[0], commandArguments[1:]

	cmd := exec.Command(head, commandArguments...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		glog.Errorf("Error during transcoding: %s\nError: %s", commandString, err.Error())
		return err
	}

	return nil
}

// TranscodeToStandard transcodes a video file to Standard preset
func TranscodeToStandard(videoName string, filename string, folderPath string) {
	glog.Infof("Transcoding to Standard preset: %s\n", videoName)

	// Standard web video
	ffmpegCommand480Pass2 := fmt.Sprintf("ffmpeg -y -i %s/%s -codec:v libx264 -profile:v high -level 4.0 -preset slow -vf scale=640:480 -threads 0 -codec:a libfdk_aac -f mp4 %s/%s_480p.mp4", folderPath, filename, folderPath, videoName)

	ExecuteFfmpegCLI(ffmpegCommand480Pass2)
}

// TranscodeToMobile transcodes video file to Mobile preset
func TranscodeToMobile(videoName string, filename string, folderPath string) {
	glog.Infof("Transcoding to Mobile preset: %s\n", videoName)

	// 360p video for older mobile phones
	ffmpegCommand360Pass2 := fmt.Sprintf("ffmpeg -y -i %s/%s -codec:v libx264 -profile:v baseline -level 3.1 -preset slow -vf scale=640:360 -threads 0 -codec:a libfdk_aac -f mp4 %s/%s_360p.mp4", folderPath, filename, folderPath, videoName)

	ExecuteFfmpegCLI(ffmpegCommand360Pass2)
}

// TranscodeToHighSD transcodes video file to HighSD preset
func TranscodeToHighSD(videoName string, filename string, folderPath string) {
	glog.Infof("Transcoding to HighSD preset: %s\n", videoName)

	// High-quality SD video for archive/storage
	ffmpegCommand576Pass2 := fmt.Sprintf("ffmpeg -y -i %s/%s -codec:v libx264 -profile:v high -level 4.2 -preset slow -vf scale=1024:576 -threads 0 -codec:a libfdk_aac -f mp4 %s/%s_576p.mp4", folderPath, filename, folderPath, videoName)

	ExecuteFfmpegCLI(ffmpegCommand576Pass2)
}
