FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN set -eux
RUN apk add --no-cache --update bash coreutils
RUN chmod +x ./scripts/wait-for-it.sh

ENTRYPOINT ./scripts/wait-for-it.sh kafka:9092 calendar:8080 calendar:8081 postgres:5432 zookeeper:2181 -- go test -tags integration ./tests/integration/...
