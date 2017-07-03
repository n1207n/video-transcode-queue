docker build -t n1207n/video-transcode-queue/video_backend_api:dev video_backend_api

docker build -t n1207n/video-transcode-queue/streaming_service:dev streaming_service

docker build -t n1207n/video-transcode-queue/transcoder_service:dev transcoder_service

docker build -t n1207n/video-transcode-queue/task_queue/consumer:dev -f task_queue/client/Dockerfile-consumer task_queue/client

docker build -t n1207n/video-transcode-queue/nginx_proxy:dev nginx_proxy
