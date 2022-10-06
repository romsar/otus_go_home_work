package sender

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"golang.org/x/sync/errgroup"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/kafka"
)

type Broker interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
}

type Model interface {
	UpdateEvent(ctx context.Context, id uuid.UUID, e *calendar.Event) (*calendar.Event, error)
}

type Sender struct {
	m   Model
	b   Broker
	cfg Config
}

type Config struct {
	Threads int
}

func New(m Model, b Broker, cfg Config) Sender {
	return Sender{
		m:   m,
		b:   b,
		cfg: cfg,
	}
}

func (s Sender) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errGrp, ctx := errgroup.WithContext(ctx)

	for i := 0; i < s.cfg.Threads; i++ {
		errGrp.Go(func() error {
			defer cancel()

			for {
				select {
				case <-ctx.Done():
					return nil
				default:
				}

				msg, err := s.b.ReadMessage(ctx)
				if err != nil {
					return err
				}

				event := new(calendar.Event)
				if err := json.Unmarshal(msg.Value, event); err != nil {
					return err
				}

				if err := sendNotification(event); err != nil {
					return err
				}

				event.IsNotified = true
				if _, err := s.m.UpdateEvent(ctx, event.ID, event); err != nil {
					return err
				}
			}
		})
	}

	if err := errGrp.Wait(); err != nil {
		return err
	}

	return nil
}

func sendNotification(e *calendar.Event) error {
	fmt.Printf("Привет, %s!\n", e.UserID)
	fmt.Printf("В %s начнется событие: %s\n", e.StartAt.Format("15:04"), e.Title)
	fmt.Printf("Описание: %s\n", e.Description)

	return nil
}
