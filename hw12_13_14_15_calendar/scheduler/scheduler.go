package scheduler

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
)

type Broker interface {
	SendEventToQueue(ctx context.Context, events ...*calendar.Event) error
}

type Repository interface {
	DeleteEvent(ctx context.Context, ids ...uuid.UUID) error
	FindEvents(ctx context.Context, filter calendar.EventFilter) ([]*calendar.Event, error)
}

type Scheduler struct {
	r   Repository
	b   Broker
	cfg Config
}

type Config struct {
	Interval        time.Duration
	EventLifeInDays uint
}

func New(r Repository, b Broker, cfg Config) Scheduler {
	return Scheduler{
		r:   r,
		b:   b,
		cfg: cfg,
	}
}

func (s Scheduler) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ticker := time.NewTicker(s.cfg.Interval)
	defer ticker.Stop()

	senderCh := make(chan struct{}, 1)
	deleteCh := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}

			senderCh <- struct{}{}
			deleteCh <- struct{}{}
		}
	}()

	errGrp, ctx := errgroup.WithContext(ctx)

	// Отправка сообщений в Sender.
	errGrp.Go(func() error {
		defer cancel()

		fn := func() error {
			events, err := s.r.FindEvents(ctx, calendar.EventFilter{
				NotNotified: true,
				NotifyTime:  true,
			})
			if err != nil {
				return err
			}

			if len(events) == 0 {
				return nil
			}

			if err := s.b.SendEventToQueue(ctx, events...); err != nil {
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
			case <-senderCh:
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
			events, err := s.r.FindEvents(ctx, calendar.EventFilter{
				To: time.Now().AddDate(0, -int(s.cfg.EventLifeInDays), 0),
			})
			if err != nil {
				return err
			}

			if len(events) == 0 {
				return nil
			}

			for _, e := range events {
				if err := s.r.DeleteEvent(ctx, e.ID); err != nil {
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
			case <-deleteCh:
			}

			if err := fn(); err != nil {
				return err
			}
		}
	})

	return errGrp.Wait()
}
