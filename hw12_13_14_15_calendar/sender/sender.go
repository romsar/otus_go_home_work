package sender

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
)

type Broker interface {
	ReadEventFromQueue(ctx context.Context) (*calendar.Event, error)
}

type Repository interface {
	UpdateEvent(ctx context.Context, id uuid.UUID, e *calendar.Event) (*calendar.Event, error)
}

type Sender struct {
	r   Repository
	b   Broker
	cfg Config
}

type Config struct {
	Threads int
}

func New(r Repository, b Broker, cfg Config) Sender {
	return Sender{
		r:   r,
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
					return ctx.Err()
				default:
				}

				event, err := s.b.ReadEventFromQueue(ctx)
				if err != nil {
					return err
				}

				if err := sendNotification(event); err != nil {
					return err
				}

				event.IsNotified = true
				if _, err := s.r.UpdateEvent(ctx, event.ID, event); err != nil {
					return err
				}
			}
		})
	}

	return errGrp.Wait()
}

func sendNotification(e *calendar.Event) error {
	fmt.Printf("Привет, %s!\n", e.UserID)
	fmt.Printf("В %s начнется событие: %s\n", e.StartAt.Format("15:04"), e.Title)
	fmt.Printf("Описание: %s\n", e.Description)

	return nil
}
