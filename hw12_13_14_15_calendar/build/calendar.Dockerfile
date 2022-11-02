FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

ARG CMD_PATH

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o ./${CMD_PATH} ./cmd/${CMD_PATH}

FROM alpine

# https://stackoverflow.com/questions/34324277/how-to-pass-arg-value-to-entrypoint
ARG CMD_PATH
ENV CMD_PATH=${CMD_PATH}

WORKDIR /app

COPY --from=builder /app .

ENTRYPOINT ./${CMD_PATH}
