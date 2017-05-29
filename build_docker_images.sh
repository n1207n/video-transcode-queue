docker build -t n1207n/video-transcode-queue/video_backend_api:dev video_backend_api

docker build -t n1207n/video-transcode-queue/task_queue/consumer:dev task_queue/client/Dockerfile-consumer

docker build -t n1207n/video-transcode-queue/nginx_proxy:dev nginx_proxy
