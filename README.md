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
  * `bash create helm_releases.sh`
3. Build docker images
 - `bash build_docker_images.sh`
4. Run Kubernetes resources
 - TODO
4. Run minikube and kubectl proxy
  * `minikube start`
  * `kubectl proxy`
