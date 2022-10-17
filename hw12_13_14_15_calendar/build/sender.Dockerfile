FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o ./calendar_sender ./cmd/calendar_sender

FROM scratch

WORKDIR /app

COPY --from=builder /app .

ENTRYPOINT ["./calendar_sender"]