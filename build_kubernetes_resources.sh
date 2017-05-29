kubectl create -f kubernetes/video_api_deployment.yml
kubectl create -f kubernetes/video_api_service.yml

kubectl create -f kubernetes/queue_consumer_job.yml
