FROM golang:1.8.3-alpine

ENV GOBIN=/go/bin

RUN apk update && apk upgrade && \
    apk add --no-cache git openssh

RUN go get -u github.com/adjust/rmq
RUN go get -u github.com/golang/glog
RUN go get -u gopkg.in/redis.v3

ADD . /go/src/github.com/n1207n/video-transcode-queue/task_queue

WORKDIR /go/src/github.com/n1207n/video-transcode-queue/task_queue

RUN go build

RUN go install

ENTRYPOINT /go/bin/task_queue

EXPOSE 6379
