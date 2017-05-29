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
2. `helm init && helm repo update`
3. Install helm packages
  * `bash create helm_releases.sh`
4. Build docker images
 - `bash build_docker_images.sh`
5. Run Kubernetes resources
 - `bash build_kubernetes_resources.sh`
6. Run minikube and kubectl proxy
  * `minikube start`
  * `kubectl proxy`
