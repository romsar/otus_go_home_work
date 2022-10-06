package kafka

import (
	"context"

	kafka "github.com/segmentio/kafka-go"
)

type Reader struct {
	r *kafka.Reader
}

type ReaderConfig = kafka.ReaderConfig

func NewReader(cfg ReaderConfig) Reader {
	return Reader{
		r: kafka.NewReader(cfg),
	}
}

func (r Reader) ReadMessage(ctx context.Context) (Message, error) {
	return r.r.ReadMessage(ctx)
}

func (r Reader) Close() error {
	return r.r.Close()
}
