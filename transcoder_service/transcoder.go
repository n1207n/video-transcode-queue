package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// ExecuteCLI executes constructed command string by os.exec.Command
func ExecuteCLI(commandString string, returnOutput bool) ([]byte, error) {
	commandArguments := strings.Fields(commandString)
	head, commandArguments := commandArguments[0], commandArguments[1:]

	cmd := exec.Command(head, commandArguments...)
	outputBytes, err := cmd.Output()

	if err != nil {
		return outputBytes, err
	}

	return outputBytes, nil
}

// GetVideoDimensionInfo extracts video width and height values
func GetVideoDimensionInfo(filename string, folderPath string) (int, int, error) {
	sugaredLogger.Infof("Getting video resolution info: %s/%s\n", folderPath, filename)

	ffprobeCommand := fmt.Sprintf("ffprobe -show_streams -print_format json -v quiet %s/%s", folderPath, filename)

	width, height := -1, -1

	outputBytes, err := ExecuteCLI(ffprobeCommand, true)
	if err != nil {
		sugaredLogger.Errorf("Error during command execution: %s\nError: %s", ffprobeCommand, err.Error())

		return width, height, err
	}

	var probeData ProbeData
	err = json.Unmarshal(outputBytes, &probeData)

	if err != nil {
		sugaredLogger.Errorf("ffprobe JSON parse error: %s\n", err.Error())

		return width, height, err
	}

	for index := 0; index < len(probeData.Stream); index++ {
		stream := probeData.Stream[index]

		if stream.Width != nil {
			width = *probeData.Stream[0].Width
			height = *probeData.Stream[0].Height
			break
		}
	}

	if width == -1 {
		return width, height, errors.New("no video stream found from file")
	}

	return width, height, nil
}

// TranscodeToSD360P transcodes video file to 360P
func TranscodeToSD360P(videoName string, videoID int, filename string, folderPath string, dbConnectionInfo map[string]string, waitGroup *sync.WaitGroup) {
	sugaredLogger.Infof("Transcoding to SD 360P: %s\n", videoName)

	defer waitGroup.Done()

	transcodedFileName := fmt.Sprintf("%s/%s_360.mp4", folderPath, videoName)

	ffmpegCommand360P := fmt.Sprintf("ffmpeg -y -i %s/%s -c:a libfdk_aac -ac 2 -ab 128k -preset slow -c:v libx264 -x264opts keyint=24:min-keyint=24:no-scenecut -b:v 400k -maxrate 400k -bufsize 400k -vf scale=-1:360 -pass 1 %s", folderPath, filename, transcodedFileName)

	_, err := ExecuteCLI(ffmpegCommand360P, false)
	if err != nil {
		sugaredLogger.Errorf("Error during command execution: %s\nError: %s", ffmpegCommand360P, err.Error())
		return
	}

	sugaredLogger.Infof("Transcoded to SD 360P: %s\n", videoName)

	width, height, err := GetVideoDimensionInfo(transcodedFileName, folderPath)
	if err != nil {
		sugaredLogger.Errorf("Error from getting video dimension info: %s\n", err.Error())
		return
	}

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

	pgDb := dbConnectionInfo["pgDb"]
	pgUser := dbConnectionInfo["pgUser"]
	pgPassword := dbConnectionInfo["pgPassword"]
	pgHost := dbConnectionInfo["pgHost"]

	connection := GetDatabaseConnection(pgUser, pgPassword, pgHost, pgDb)
	CreateVideoRenderingObject(videoRendering, connection)
}

// TranscodeToSD540P transcodes video file to 540P
func TranscodeToSD540P(videoName string, videoID int, filename string, folderPath string, dbConnectionInfo map[string]string, waitGroup *sync.WaitGroup) {
	sugaredLogger.Infof("Transcoding to SD 540P: %s\n", videoName)

	defer waitGroup.Done()

	transcodedFileName := fmt.Sprintf("%s/%s_540.mp4", folderPath, videoName)

	ffmpegCommand540P := fmt.Sprintf("ffmpeg -y -i %s/%s -c:a libfdk_aac -ac 2 -ab 128k -preset slow -c:v libx264 -x264opts keyint=24:min-keyint=24:no-scenecut -b:v 800k -maxrate 800k -bufsize 500k -vf scale=-1:540 -pass 1 %s", folderPath, filename, transcodedFileName)

	_, err := ExecuteCLI(ffmpegCommand540P, false)
	if err != nil {
		sugaredLogger.Errorf("Error during command execution: %s\nError: %s", ffmpegCommand540P, err.Error())
		return
	}

	width, height, err := GetVideoDimensionInfo(transcodedFileName, folderPath)
	if err != nil {
		sugaredLogger.Errorf("Error from getting video dimension info: %s\n", err.Error())
		return
	}

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

	pgDb := dbConnectionInfo["pgDb"]
	pgUser := dbConnectionInfo["pgUser"]
	pgPassword := dbConnectionInfo["pgPassword"]
	pgHost := dbConnectionInfo["pgHost"]

	connection := GetDatabaseConnection(pgUser, pgPassword, pgHost, pgDb)
	CreateVideoRenderingObject(videoRendering, connection)
}

// TranscodeToHD720P transcodes video file to 720P
func TranscodeToHD720P(videoName string, videoID int, filename string, folderPath string, dbConnectionInfo map[string]string, waitGroup *sync.WaitGroup) {
	sugaredLogger.Infof("Transcoding to HD 720P: %s\n", videoName)

	defer waitGroup.Done()

	transcodedFileName := fmt.Sprintf("%s/%s_720.mp4", folderPath, videoName)

	ffmpegCommand720P := fmt.Sprintf("ffmpeg -y -i %s/%s -c:a libfdk_aac -ac 2 -ab 128k -preset slow -c:v libx264 -x264opts keyint=24:min-keyint=24:no-scenecut -b:v 1500k -maxrate 1500k -bufsize 1000k -vf scale=-1:720 -pass 1 %s", folderPath, filename, transcodedFileName)

	_, err := ExecuteCLI(ffmpegCommand720P, false)
	if err != nil {
		sugaredLogger.Errorf("Error during command execution: %s\nError: %s", ffmpegCommand720P, err.Error())
		return
	}

	width, height, err := GetVideoDimensionInfo(transcodedFileName, folderPath)
	if err != nil {
		sugaredLogger.Errorf("Error from getting video dimension info: %s\n", err.Error())
		return
	}

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

	pgDb := dbConnectionInfo["pgDb"]
	pgUser := dbConnectionInfo["pgUser"]
	pgPassword := dbConnectionInfo["pgPassword"]
	pgHost := dbConnectionInfo["pgHost"]

	connection := GetDatabaseConnection(pgUser, pgPassword, pgHost, pgDb)
	CreateVideoRenderingObject(videoRendering, connection)
}

// ConstructMPD creates MPD file for DASH streaming
func ConstructMPD(videoName string, videoID int, filename string, folderPath string, transcodeTargets []int, dbConnectionInfo map[string]string) {
	sugaredLogger.Infof("Constructing MPD file: %s\n", videoName)

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
		sugaredLogger.Errorf("Error during command execution: %s\nError: %s", mp4boxCommand, err.Error())
	} else {
		video := Video{
			UpdatedAt:      time.Now(),
			StreamFilePath: fmt.Sprintf("%s.mpd", filePath),
			IsReadyToServe: true,
		}

		pgDb := dbConnectionInfo["pgDb"]
		pgUser := dbConnectionInfo["pgUser"]
		pgPassword := dbConnectionInfo["pgPassword"]
		pgHost := dbConnectionInfo["pgHost"]

		connection := GetDatabaseConnection(pgUser, pgPassword, pgHost, pgDb)
		UpdateVideoObject(video, connection)
	}
}
