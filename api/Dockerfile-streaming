FROM golang:1.8.3-alpine

ENV GOBIN=/go/bin

RUN apk update && apk upgrade && \
    apk add --no-cache git openssh

RUN go get -u github.com/satori/go.uuid
RUN go get -u github.com/joho/godotenv
RUN go get -u go.uber.org/zap
RUN go get -u github.com/gin-gonic/gin
RUN go get -u github.com/gin-contrib/static
RUN go get -u github.com/gin-gonic/autotls
RUN go get -u github.com/gin-contrib/cors
RUN go get -u github.com/jinzhu/gorm
RUN go get -u github.com/jinzhu/gorm/dialects/postgres
RUN go get -u gopkg.in/redis.v3
RUN go get -u github.com/adjust/rmq

ADD common /go/src/github.com/n1207n/video-transcode-queue/api/common
ADD streaming /go/src/github.com/n1207n/video-transcode-queue/api/streaming

WORKDIR /go/src/github.com/n1207n/video-transcode-queue/api/streaming

RUN go build
RUN go install

ENTRYPOINT /go/bin/streaming

EXPOSE 8880
