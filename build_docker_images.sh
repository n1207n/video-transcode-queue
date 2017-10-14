docker build -f api/Dockerfile-video-backend -t n1207n/video-transcode-queue/video_backend_api:dev api

docker build -f api/Dockerfile-streaming -t n1207n/video-transcode-queue/streaming_service:dev api

docker build -f api/Dockerfile-transcoder -t n1207n/video-transcode-queue/transcoder_service:dev api

docker build -t n1207n/video-transcode-queue/task_queue/consumer:dev -f task_queue/client/Dockerfile-consumer task_queue/client
