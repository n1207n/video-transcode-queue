FROM golang:1.8.3-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git openssh

RUN go get -u github.com/kataras/iris/iris
RUN go get -u github.com/go-pg/pg
RUN go get -u github.com/satori/go.uuid
RUN go get -u github.com/jmoiron/sqlx
RUN go get -u github.com/jmoiron/sqlx/types
RUN go get -u github.com/joho/godotenv

ADD . /code/
WORKDIR /code
CMD ["go", "run", "main.go"]
