package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Writer struct {
	w *kafka.Writer
}

type WriterConfig = kafka.Writer

func NewWriter(cfg *WriterConfig) Writer {
	return Writer{
		w: cfg,
	}
}

func (w Writer) Close() error {
	return w.w.Close()
}

func (w Writer) WriteMessages(ctx context.Context, msgs ...Message) error {
	return w.w.WriteMessages(ctx, msgs...)
}
