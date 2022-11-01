package kafka

import (
	"context"

	json "github.com/json-iterator/go"
	kafka "github.com/segmentio/kafka-go"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
)

type Writer struct {
	w *kafka.Writer
}

type WriterConfig struct {
	Brokers []string
	Topic   string
}

func NewWriter(cfg *WriterConfig) Writer {
	return Writer{
		w: &kafka.Writer{
			Addr:  kafka.TCP(cfg.Brokers...),
			Topic: cfg.Topic,
		},
	}
}

func (w Writer) Close() error {
	return w.w.Close()
}

func (w Writer) SendEventToQueue(ctx context.Context, events ...*calendar.Event) error {
	messages := make([]kafka.Message, 0, len(events))
	for _, e := range events {
		bs, err := json.Marshal(e)
		if err != nil {
			return err
		}

		messages = append(messages, kafka.Message{
			Value: bs,
		})
	}

	return w.w.WriteMessages(ctx, messages...)
}
