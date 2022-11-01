package kafka

import (
	"context"

	json "github.com/json-iterator/go"
	kafka "github.com/segmentio/kafka-go"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
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

func (r Reader) Close() error {
	return r.r.Close()
}

func (r Reader) ReadEventFromQueue(ctx context.Context) (*calendar.Event, error) {
	msg, err := r.r.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}

	event := new(calendar.Event)
	if err := json.Unmarshal(msg.Value, event); err != nil {
		return nil, err
	}

	return event, nil
}
