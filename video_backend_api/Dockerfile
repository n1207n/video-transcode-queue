FROM golang:1.8.3-alpine

ENV GOBIN=/go/bin

RUN apk update && apk upgrade && \
    apk add --no-cache git openssh

RUN go get -u github.com/satori/go.uuid
RUN go get -u github.com/joho/godotenv
RUN go get -u github.com/golang/glog
RUN go get -u gopkg.in/gin-gonic/gin.v1
RUN go get -u github.com/jinzhu/gorm
RUN go get -u github.com/jinzhu/gorm/dialects/postgres
RUN go get -u gopkg.in/redis.v3
RUN go get -u github.com/adjust/rmq

ADD . /go/src/github.com/n1207n/video-transcode-queue/video_backend_api

RUN go install /go/src/github.com/n1207n/video-transcode-queue/video_backend_api/*.go

ENTRYPOINT /go/bin/api_server

EXPOSE 8080
