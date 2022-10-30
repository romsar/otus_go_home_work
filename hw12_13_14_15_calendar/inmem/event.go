package inmem

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
)

var TimeNowFunc func() time.Time

func init() {
	if TimeNowFunc == nil {
		TimeNowFunc = time.Now
	}
}

// CreateEvent создает событие.
func (repo *Repository) CreateEvent(ctx context.Context, e *calendar.Event) (*calendar.Event, error) {
	if err := repo.checkDateBusy(e); err != nil {
		return nil, errors.Wrap(err, "create event")
	}

	e.ID = uuid.New()

	repo.eventMu.Lock()
	defer repo.eventMu.Unlock()

	repo.events[e.ID] = e

	return e, nil
}

// UpdateEvent обновляет событие.
func (repo *Repository) UpdateEvent(ctx context.Context, id uuid.UUID, e *calendar.Event) (*calendar.Event, error) {
	if _, err := repo.findEventByID(id); err != nil {
		return nil, errors.Wrap(err, "update event")
	}

	if err := repo.checkDateBusy(e, id); err != nil {
		return nil, errors.Wrap(err, "update event")
	}

	repo.eventMu.Lock()
	defer repo.eventMu.Unlock()

	repo.events[id] = e
	repo.events[id].ID = id

	return e, nil
}

// DeleteEvent удаляет событие.
func (repo *Repository) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	if _, err := repo.findEventByID(id); err != nil {
		return errors.Wrap(err, "delete event")
	}

	repo.eventMu.Lock()
	defer repo.eventMu.Unlock()

	delete(repo.events, id)

	return nil
}

// FindEvents находит события по критериям.
func (repo *Repository) FindEvents(ctx context.Context, filter calendar.EventFilter) ([]*calendar.Event, error) {
	repo.eventMu.Lock()
	defer repo.eventMu.Unlock()

	res := make([]*calendar.Event, 0)

	for _, e := range repo.events {
		if !passFilter(e, filter) {
			continue
		}

		res = append(res, e)
	}

	return res, nil
}

// FindEventByID находит событие по ID.
func (repo *Repository) FindEventByID(ctx context.Context, id uuid.UUID) (*calendar.Event, error) {
	event, err := repo.findEventByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "find event")
	}

	return event, nil
}

// findEventByID находит событие по ID.
func (repo *Repository) findEventByID(id uuid.UUID) (*calendar.Event, error) {
	repo.eventMu.Lock()
	defer repo.eventMu.Unlock()

	event, exists := repo.events[id]
	if !exists {
		return nil, calendar.ErrNotFound
	}

	return event, nil
}

// passFilter проверяет событие на удовлетворенность условиям фильтра.
func passFilter(e *calendar.Event, filter calendar.EventFilter) bool {
	if filter.UserID != uuid.Nil && e.UserID != filter.UserID {
		return false
	}

	if filter.NotNotified && e.IsNotified {
		return false
	}

	if filter.NotifyTime {
		now := TimeNowFunc()

		if e.StartAt.Before(now) {
			return false
		}

		notifyAt := e.StartAt.Add(-time.Duration(int64(e.NotificationDuration)) * time.Minute)

		if now.Before(notifyAt) {
			return false
		}
	}

	if !filter.From.IsZero() && e.StartAt.Before(filter.From) {
		return false
	}

	if !filter.To.IsZero() && e.EndAt.After(filter.To) {
		return false
	}

	return true
}

// checkDateBusy проверка на свободное время.
// Если время занято, то вернет ошибку calendar.ErrDateBusy.
func (repo *Repository) checkDateBusy(event *calendar.Event, ID ...uuid.UUID) error {
	var ignore uuid.UUID
	if len(ID) > 0 {
		ignore = ID[0]
	}
	
	repo.eventMu.Lock()
	defer repo.eventMu.Unlock()

	for _, e := range repo.events {
		if e.ID == ignore {
			continue
		}

		if e.UserID != event.UserID {
			continue
		}

		if (e.StartAt.After(event.StartAt) || e.StartAt.Equal(event.StartAt)) &&
			e.StartAt.Before(event.EndAt) {
			return calendar.ErrDateBusy
		}

		if e.EndAt.After(event.StartAt) && (e.EndAt.Before(event.EndAt) || e.EndAt.Equal(event.EndAt)) {
			return calendar.ErrDateBusy
		}
	}

	return nil
}
