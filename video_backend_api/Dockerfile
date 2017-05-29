FROM golang:1.8.3-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git openssh

RUN go get -u github.com/kataras/iris/iris
RUN go get -u github.com/go-pg/pg
RUN go get -u github.com/satori/go.uuid
RUN go get -u github.com/jmoiron/sqlx
RUN go get -u github.com/jmoiron/sqlx/types
RUN go get -u github.com/joho/godotenv
RUN go get -u github.com/golang/glog

ADD . /go/src/github.com/n1207n/video-transcode-queue/video_backend_api

RUN go install github.com/n1207n/video-transcode-queue/video_backend_api/api_server

WORKDIR /go/src/github.com/n1207n/video-transcode-queue/video_backend_api

ENTRYPOINT /go/bin/api_server

EXPOSE 8000
