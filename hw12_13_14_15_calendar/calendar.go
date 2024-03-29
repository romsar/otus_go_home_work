package calendar

import (
	"context"

	"github.com/google/uuid"
)

// Repository декларирует контракт репозитория.
type Repository interface {
	// CreateEvent создать событие.
	CreateEvent(ctx context.Context, e *Event) (*Event, error)

	// UpdateEvent обновить событие.
	UpdateEvent(ctx context.Context, id uuid.UUID, e *Event) (*Event, error)

	// DeleteEvent удалить событие.
	DeleteEvent(ctx context.Context, ids ...uuid.UUID) error

	// FindEvents найти множество событий.
	FindEvents(ctx context.Context, filter EventFilter) ([]*Event, error)

	// FindEventByID найти событие по его идентификатору.
	FindEventByID(ctx context.Context, id uuid.UUID) (*Event, error)
}
