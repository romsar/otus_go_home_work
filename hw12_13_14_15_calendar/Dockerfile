FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY . .

RUN echo
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app/calendar ./cmd/calendar

FROM alpine

COPY --from=builder /app .

EXPOSE 8080
CMD ["./calendar"]