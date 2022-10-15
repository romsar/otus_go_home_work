FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY .  .

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

RUN go build -o ./calendar_scheduler ./cmd/calendar_scheduler

FROM scratch

WORKDIR /app

COPY --from=builder /app .

CMD ["./calendar_scheduler"]