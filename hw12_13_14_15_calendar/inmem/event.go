package inmem

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
)

// CreateEvent создает событие.
func (repo *Repository) CreateEvent(ctx context.Context, e *calendar.Event) (*calendar.Event, error) {
	e.ID = uuid.New()

	repo.events[e.ID] = e

	return e, nil
}

// UpdateEvent обновляет событие.
func (repo *Repository) UpdateEvent(ctx context.Context, id uuid.UUID, e *calendar.Event) (*calendar.Event, error) {
	if _, err := repo.findEventByID(id); err != nil {
		return nil, errors.Wrap(err, "update event")
	}

	repo.events[id] = e
	repo.events[id].ID = id

	return e, nil
}

// DeleteEvent удаляет событие.
func (repo *Repository) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	if _, err := repo.findEventByID(id); err != nil {
		return errors.Wrap(err, "delete event")
	}

	delete(repo.events, id)

	return nil
}

// FindEvents находит события по критериям.
func (repo *Repository) FindEvents(ctx context.Context, filter calendar.EventFilter) ([]*calendar.Event, error) {
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

	if !filter.From.IsZero() && filter.From.Before(e.StartAt) {
		return false
	}

	if !filter.To.IsZero() && filter.To.After(e.StartAt.Add(e.Duration)) {
		return false
	}

	return true
}
