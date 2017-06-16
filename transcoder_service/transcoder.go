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

	ffmpegCommand480Pass1 := fmt.Sprintf("ffmpeg -y -i %s/%s -codec:v libx264 -profile:v high -preset slow -b:v 500k -maxrate 500k -bufsize 1000k -vf scale=-2:480 -threads 0 -pass 1 -an -f mp4 /dev/null", folderPath, filename)

	ExecuteFfmpegCLI(ffmpegCommand480Pass1)

	// Standard web video (480p at 500kbit/s)
	ffmpegCommand480Pass2 := fmt.Sprintf("ffmpeg -y -i %s/%s -codec:v libx264 -profile:v high -preset slow -b:v 500k -maxrate 500k -bufsize 1000k -vf scale=-2:480 -threads 0 -pass 2 -codec:a libfdk_aac -b:a 128k -f mp4 %s/%s_480p.mp4", folderPath, filename, folderPath, videoName)

	ExecuteFfmpegCLI(ffmpegCommand480Pass2)
}

// TranscodeToMobile transcodes video file to Mobile preset
func TranscodeToMobile(videoName string, filename string, folderPath string) {
	glog.Infof("Transcoding to Mobile preset: %s\n", videoName)

	ffmpegCommand360Pass1 := fmt.Sprintf("ffmpeg -y -i %s/%s -codec:v libx264 -profile:v baseline -level 3.1 -preset slow -b:v 250k -maxrate 250k -bufsize 500k -vf scale=-2:360 -threads 0 -pass 1 -an -f mp4 /dev/null", folderPath, filename)

	ExecuteFfmpegCLI(ffmpegCommand360Pass1)

	// 360p video for older mobile phones (360p at 250kbit/s in baseline profile)
	ffmpegCommand360Pass2 := fmt.Sprintf("ffmpeg -y -i %s/%s -codec:v libx264 -profile:v baseline -level 3.1 -preset slow -b:v 250k -maxrate 250k -bufsize 500k -vf scale=-2:360 -threads 0 -pass 2 -codec:a libfdk_aac -b:a 128k -f mp4 %s/%s_360p.mp4", folderPath, filename, folderPath, videoName)

	ExecuteFfmpegCLI(ffmpegCommand360Pass2)
}

// TranscodeToHighSD transcodes video file to HighSD preset
func TranscodeToHighSD(videoName string, filename string, folderPath string) {
	glog.Infof("Transcoding to HighSD preset: %s\n", videoName)

	ffmpegCommand576Pass1 := fmt.Sprintf("ffmpeg -y -i %s/%s -codec:v libx264 -profile:v high -preset slow -b:v 1000k -vf scale=-2:576 -threads 0 -pass 1 -an -f mp4 /dev/null", folderPath, filename)

	ExecuteFfmpegCLI(ffmpegCommand576Pass1)

	// High-quality SD video for archive/storage (PAL at 1Mbit/s in high profile):
	ffmpegCommand576Pass2 := fmt.Sprintf("ffmpeg -y -i %s/%s -codec:v libx264 -profile:v high -preset slow -b:v 1000k -vf scale=-2:576 -threads 0 -pass 2 -codec:a libfdk_aac -b:a 196k -f mp4 %s/%s_576p.mp4", folderPath, filename, folderPath, videoName)

	ExecuteFfmpegCLI(ffmpegCommand576Pass2)
}
