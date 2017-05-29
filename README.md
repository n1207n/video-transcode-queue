# Video-transcode-queue
A sample infrastructure for processing video upload & transcoding.
* React web UI
* REST API in Golang
* PostgreSQL
* Nginx proxy
* Redis task queue storage
* Redis task consumer in Golang
* Redis task producer in Golang
* Video transcoder in Golang

## Powered by Kubernetes Helm packages

### How to run locally
1. Install Docker, Kubernetes, Minikube, and helm package manager
2. Install redis message broker via helm
  * `helm install -f task_queue/queue_storage_config.yml --name=queue-storage --namespace dev-cluster-1 stable/redis`
3. Run minikube and kubectl proxy
  * `minikube start`
  * `kubectl proxy`
