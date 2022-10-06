package scheduler

import (
	"context"
	"time"

	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"golang.org/x/sync/errgroup"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/kafka"
)

type Broker interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

type Model interface {
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	FindEvents(ctx context.Context, filter calendar.EventFilter) ([]*calendar.Event, error)
}

type Scheduler struct {
	m   Model
	b   Broker
	cfg Config
}

type Config struct {
	Interval        time.Duration
	EventLifeInDays uint
}

func New(m Model, b Broker, cfg Config) Scheduler {
	return Scheduler{
		m:   m,
		b:   b,
		cfg: cfg,
	}
}

func (s Scheduler) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.cfg.Interval)
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errGrp, ctx := errgroup.WithContext(ctx)

	// Отправка сообщений в Sender.
	errGrp.Go(func() error {
		defer cancel()

		fn := func() error {
			events, err := s.m.FindEvents(ctx, calendar.EventFilter{
				NotNotified: true,
				NotifyTime:  true,
			})
			if err != nil {
				return err
			}

			if len(events) == 0 {
				return nil
			}

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

			if err := s.b.WriteMessages(ctx, messages...); err != nil {
				return err
			}

			return nil
		}

		if err := fn(); err != nil {
			return err
		}

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
			}

			if err := fn(); err != nil {
				return err
			}
		}
	})

	// Удаление старых событий.
	errGrp.Go(func() error {
		defer cancel()

		fn := func() error {
			events, err := s.m.FindEvents(ctx, calendar.EventFilter{
				To: time.Now().AddDate(0, -int(s.cfg.EventLifeInDays), 0),
			})
			if err != nil {
				return err
			}

			if len(events) == 0 {
				return nil
			}

			for _, e := range events {
				if err := s.m.DeleteEvent(ctx, e.ID); err != nil {
					return err
				}
			}

			return nil
		}

		if err := fn(); err != nil {
			return err
		}

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
			}

			if err := fn(); err != nil {
				return err
			}
		}
	})

	if err := errGrp.Wait(); err != nil {
		return err
	}

	return nil
}
