package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
)

// ExecuteCLI executes constructed command string by os.exec.Command
func ExecuteCLI(commandString string, returnOutput bool) (*bytes.Buffer, error) {
	commandArguments := strings.Fields(commandString)
	head, commandArguments := commandArguments[0], commandArguments[1:]

	cmd := exec.Command(head, commandArguments...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	outputBytes := new(bytes.Buffer)

	if err := cmd.Run(); err != nil {
		return outputBytes, err
	}

	if returnOutput != false {
		io.Copy(cmd.Stdout, outputBytes)
	}

	return outputBytes, nil
}

// GetVideoDimensionInfo extracts video width and height values
func GetVideoDimensionInfo(filename string, folderPath string) (int, int) {
	glog.Infof("Getting video resolution info: %s/%s\n", folderPath, filename)

	// ffprobe should return output as the below:
	// width=1280
	// height=720
	ffprobeCommand := fmt.Sprintf("ffprobe -v error -show_entries stream=width,height -of default=noprint_wrappers=1 %s/%s", folderPath, filename)

	outputBytes, err := ExecuteCLI(ffprobeCommand, true)
	if err != nil {
		glog.Errorf("Error during command execution: %s\nError: %s", ffprobeCommand, err.Error())
	}

	outputString := outputBytes.String()
	outputTokens := strings.Split(outputString, "\n")

	width, _ := strconv.Atoi(strings.Split(outputTokens[0], "=")[1])

	height, _ := strconv.Atoi(strings.Split(outputTokens[1], "=")[1])

	return width, height
}

// TranscodeToSD360P transcodes video file to 360P
func TranscodeToSD360P(videoName string, videoID int, filename string, folderPath string, waitGroup *sync.WaitGroup) {
	glog.Infof("Transcoding to SD 360P: %s\n", videoName)
	waitGroup.Add(1)

	ffmpegCommand360P := fmt.Sprintf("ffmpeg -y -i %s/%s -c:a libfdk_aac -ac 2 -ab 128k -preset slow -c:v libx264 -x264opts keyint=24:min-keyint=24:no-scenecut -b:v 400k -maxrate 400k -bufsize 400k -vf scale=-1:360 -pass 1 %s/%s_360.mp4", folderPath, filename, folderPath, videoName)

	_, err := ExecuteCLI(ffmpegCommand360P, false)
	if err != nil {
		glog.Errorf("Error during command execution: %s\nError: %s", ffmpegCommand360P, err.Error())
	}

	width, height := GetVideoDimensionInfo(filename, folderPath)

	videoRendering := VideoRendering{
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		RenderingTitle: fmt.Sprintf("%s_360", videoName),
		FilePath:       fmt.Sprintf("%s/%s_360.mp4", folderPath, videoName),
		URL:            fmt.Sprintf("%s/%s_360.mp4", folderPath, videoName),
		Width:          uint(width),
		Height:         uint(height),
		VideoID:        uint(videoID),
	}

	connection := GetDatabaseConnection(pgUser, pgPassword, pgHost, pgDb)
	CreateVideoRenderingObject(videoRendering, connection)

	waitGroup.Done()
}

// TranscodeToSD540P transcodes video file to 540P
func TranscodeToSD540P(videoName string, videoID int, filename string, folderPath string, waitGroup *sync.WaitGroup) {
	glog.Infof("Transcoding to SD 540P: %s\n", videoName)
	waitGroup.Add(1)

	ffmpegCommand540P := fmt.Sprintf("ffmpeg -y -i %s/%s -c:a libfdk_aac -ac 2 -ab 128k -preset slow -c:v libx264 -x264opts keyint=24:min-keyint=24:no-scenecut -b:v 800k -maxrate 800k -bufsize 500k -vf scale=-1:540 -pass 1 %s/%s_540.mp4", folderPath, filename, folderPath, videoName)

	_, err := ExecuteCLI(ffmpegCommand540P, false)
	if err != nil {
		glog.Errorf("Error during command execution: %s\nError: %s", ffmpegCommand540P, err.Error())
	}

	width, height := GetVideoDimensionInfo(filename, folderPath)

	videoRendering := VideoRendering{
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		RenderingTitle: fmt.Sprintf("%s_540", videoName),
		FilePath:       fmt.Sprintf("%s/%s_540.mp4", folderPath, videoName),
		URL:            fmt.Sprintf("%s/%s_540.mp4", folderPath, videoName),
		Width:          uint(width),
		Height:         uint(height),
		VideoID:        uint(videoID),
	}

	connection := GetDatabaseConnection(pgUser, pgPassword, pgHost, pgDb)
	CreateVideoRenderingObject(videoRendering, connection)

	waitGroup.Done()
}

// TranscodeToHD720P transcodes video file to 720P
func TranscodeToHD720P(videoName string, videoID int, filename string, folderPath string, waitGroup *sync.WaitGroup) {
	glog.Infof("Transcoding to HD 720P: %s\n", videoName)
	waitGroup.Add(1)

	ffmpegCommand720P := fmt.Sprintf("ffmpeg -y -i %s/%s -c:a libfdk_aac -ac 2 -ab 128k -preset slow -c:v libx264 -x264opts keyint=24:min-keyint=24:no-scenecut -b:v 1500k -maxrate 1500k -bufsize 1000k -vf scale=-1:720 -pass 1 %s/%s_720.mp4", folderPath, filename, folderPath, videoName)

	_, err := ExecuteCLI(ffmpegCommand720P, false)
	if err != nil {
		glog.Errorf("Error during command execution: %s\nError: %s", ffmpegCommand720P, err.Error())
	}

	width, height := GetVideoDimensionInfo(filename, folderPath)

	videoRendering := VideoRendering{
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		RenderingTitle: fmt.Sprintf("%s_720", videoName),
		FilePath:       fmt.Sprintf("%s/%s_720.mp4", folderPath, videoName),
		URL:            fmt.Sprintf("%s/%s_720.mp4", folderPath, videoName),
		Width:          uint(width),
		Height:         uint(height),
		VideoID:        uint(videoID),
	}

	connection := GetDatabaseConnection(pgUser, pgPassword, pgHost, pgDb)
	CreateVideoRenderingObject(videoRendering, connection)

	waitGroup.Done()
}

// ConstructMPD creates MPD file for DASH streaming
func ConstructMPD(videoName string, videoID int, filename string, folderPath string, transcodeTargets []int) {
	glog.Infof("Constructing MPD file: %s\n", videoName)

	filePath := fmt.Sprintf("%s/%s", folderPath, videoName)

	mp4boxCommand := fmt.Sprintf("MP4Box -dash 3000 -frag 3000 -rap -profile dashavc264:onDemand -out %s.mpd", filePath)

	// Appending video streams for each transcoded size
	for resize := range transcodeTargets {
		mp4boxCommand += fmt.Sprintf(" %s_%d.mp4#video", filePath, resize)
	}

	// Appending audio streams for each transcoded size
	for resize := range transcodeTargets {
		mp4boxCommand += fmt.Sprintf(" %s_%d.mp4#audio", filePath, resize)
	}

	_, err := ExecuteCLI(mp4boxCommand, false)
	if err != nil {
		glog.Errorf("Error during command execution: %s\nError: %s", mp4boxCommand, err.Error())
	}

}
