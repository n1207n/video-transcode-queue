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
2. Install ffmpeg && codecs (If you want to troubleshoot ffmpeg commands outside of docker)
 * `brew install ffmpeg --with-chromaprint --with-fdk-aac --with-fontconfig --with-freetype --with-frei0r --with-game-music-emu --with-libass --with-libbluray --with-libbs2b --with-libcaca --with-libebur128 --with-libgsm --with-libmodplug --with-libsoxr --with-libssh --with-libvidstab --with-libvorbis --with-libvpx --with-opencore-amr --with-openh264 --with-openjpeg --with-openssl --with-opus --with-rtmpdump --with-rubberband --with-schroedinger --with-sdl2 --with-snappy --with-speex --with-tesseract --with-theora --with-tools --with-two-lame --with-wavpack --with-webp --with-x265 --with-xz --with-zeromq --with-zimg`
2. Expose FFMPEG C libraries (If you want to troubleshoot ffmpeg commands outside of docker)
  * `export FFMPEG_ROOT=export FFMPEG_ROOT=/usr/local/Cellar/ffmpeg/3.3.1
export CGO_LDFLAGS="-L$FFMPEG_ROOT/lib/ -lavcodec -lavformat -lavutil -lswscale -lswresample -lavdevice -lavfilter"
export CGO_CFLAGS="-I$FFMPEG_ROOT/include"
export LD_LIBRARY_PATH=$HOME/ffmpeg/lib`
2. `helm init && helm repo update`
3. Install helm packages
  * `bash build_helm_packages.sh`
4. Build docker images
 - `eval $(minikube docker-env)`
 - `bash build_docker_images.sh`
5. Run Kubernetes resources
 - `bash build_kubernetes_resources.sh`
6. Run minikube and kubectl proxy
  * `minikube start`
  * `kubectl proxy`
7, Access minikube external url
  * `minikube service video-api --url` or `minikube service streaming-api --url`
