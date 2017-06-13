kubectl create -f kubernetes/secrets/redis-queue-info.yml

kubectl create -f kubernetes/video-upload-minikube-persistent-volume.yml
kubectl create -f kubernetes/video-upload-minikube-persistent-volume-claim.yml

kubectl create -f kubernetes/video-api-deployment.yml
# kubectl create -f kubernetes/video-api-service.yml
kubectl expose deployment video-api --type=NodePort

kubectl create -f kubernetes/transcoder-api-deployment.yml
kubectl create -f kubernetes/transcoder-api-service.yml

kubectl create -f kubernetes/nginx-deployment.yml
kubectl expose deployment nginx-proxy --type=NodePort

kubectl create -f kubernetes/queue-consumer-job.yml
