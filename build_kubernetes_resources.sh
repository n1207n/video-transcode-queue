kubectl create -f kubernetes/minikube/secrets/redis-queue-info.yml

kubectl create -f kubernetes/minikube/video-upload-minikube-persistent-volume.yml
kubectl create -f kubernetes/minikube/video-upload-minikube-persistent-volume-claim.yml

kubectl create -f kubernetes/minikube/video-api-deployment.yml
# kubectl create -f kubernetes/minikube/video-api-service.yml
kubectl expose deployment video-api --type=NodePort

kubectl create -f kubernetes/minikube/streaming-api-deployment.yml
kubectl expose deployment streaming-api --type=NodePort
# kubectl create -f kubernetes/minikube/streaming-api-service.yml

kubectl create -f kubernetes/minikube/transcoder-api-deployment.yml
kubectl create -f kubernetes/minikube/transcoder-api-service.yml

kubectl create -f kubernetes/minikube/queue-consumer-job.yml
