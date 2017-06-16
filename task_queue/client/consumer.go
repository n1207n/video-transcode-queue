package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/adjust/rmq"
	"github.com/golang/glog"
	"gopkg.in/redis.v3"
)

var (
	redisURL, redisPort, redisPassword, redisTopic string
	redisProtocol                                  = "tcp"
	redisNetworkTag                                = "transcode_task_consume"
	transcodeServiceHost, transcodeServicePort     string

	// TODO: Make below variables as a CLI argument
	queueFetchInterval            = 10
	queueFetchIntervalMeasurement = time.Second
)

func main() {
	loadEnvironmentVariables()

	taskQueue := openTaskQueue()
	glog.Infof("Queue accessed: %s\n", redisTopic)

	taskQueue.StartConsuming(queueFetchInterval, queueFetchIntervalMeasurement)
	glog.Infoln("Queue consumption started...")

	taskConsumer := &TaskConsumer{}
	taskQueue.AddConsumer("Task consumer 1", taskConsumer)

	select {}
}

// TaskConsumer represents the Redis topic consumer
type TaskConsumer struct {
	name         string
	count        int
	lastAccessed time.Time
}

// TranscodeRequest represents a JSON POST data for video-transcode API
type TranscodeRequest struct {
	Path    string `json:"path"`
	VideoID string `json:"video_id"`
}

// Consume method implements TaskConsumer struct
// to be registered on Queue.
// It handles actual data handling from Queue
func (tc *TaskConsumer) Consume(delivery rmq.Delivery) {
	var task Task

	tc.count++

	if err := json.Unmarshal([]byte(delivery.Payload()), &task); err != nil {
		glog.Errorf("Failed to read task message: %s\n", err)
		delivery.Reject()
		return
	}

	glog.Infof("Processed task message: Transcoding %s\n", task.FilePath)

	// TODO: Call Go subroutine to call go binding of ffmpeg
	transcodeRequest := &TranscodeRequest{
		Path:    task.FilePath,
		VideoID: task.ID,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(transcodeRequest)

	url := fmt.Sprintf("http://%s:%s/api/v1/video-transcode", transcodeServiceHost, transcodeServicePort)

	request, err := http.NewRequest("POST", url, b)
	if err != nil {
		glog.Warningf("Failed to trigger transcode API: %s\n", err)
		delivery.Reject()
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		glog.Warningf("Unsuccessful transcode request: %s\n", err)
		delivery.Reject()
		return
	}

	responseBuffer := new(bytes.Buffer)
	io.Copy(responseBuffer, response.Body)

	glog.Infof("Successful transcode request: %s\n", responseBuffer)
	delivery.Ack()
}

// loadEnvironmentVariables loads Redis
// information from environment variables
func loadEnvironmentVariables() {
	redisURL = os.Getenv("REDIS_URL")
	if len(redisURL) == 0 {
		panic("No REDIS_URL environment variable")
	}

	redisPort = os.Getenv("REDIS_PORT")
	if len(redisPort) == 0 {
		panic("No REDIS_PORT environment variable")
	}

	redisTopic = os.Getenv("REDIS_TOPIC")
	if len(redisTopic) == 0 {
		panic("No REDIS_TOPIC environment variable")
	}

	redisPassword = os.Getenv("REDIS_PASSWORD")
	if len(redisPassword) == 0 {
		panic("No REDIS_PASSWORD environment variable")
	}

	redisTopic = os.Getenv("REDIS_TOPIC")
	if len(redisTopic) == 0 {
		panic("No REDIS_TOPIC environment variable")
	}

	transcodeServiceHost = os.Getenv("TRANSCODER_API_SERVICE_HOST")
	if len(transcodeServiceHost) == 0 {
		panic("No TRANSCODER_API_SERVICE_HOST environment variable")
	}

	transcodeServicePort = os.Getenv("TRANSCODER_API_SERVICE_PORT")
	if len(transcodeServicePort) == 0 {
		panic("No TRANSCODER_API_SERVICE_PORT environment variable")
	}
}

// openTaskQueue connects to redis and return a Queue interface
func openTaskQueue() rmq.Queue {
	redisClient := redis.NewClient(&redis.Options{
		Network:  redisProtocol,
		Addr:     fmt.Sprintf("%s:%s", redisURL, redisPort),
		DB:       int64(1),
		Password: redisPassword,
	})

	connection := rmq.OpenConnectionWithRedisClient(redisNetworkTag, redisClient)

	glog.Infof("Connected to Redis task queue: %s\n", connection.Name)

	return connection.OpenQueue(redisTopic)
}
