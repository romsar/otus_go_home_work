package calendar

import (
	"time"

	"github.com/google/uuid"
)

// Event (событие) - основная сущность приложения.
type Event struct {
	// ID уникальный идентификатор события.
	ID uuid.UUID `db:"id"`

	// Title заголовок события.
	Title string `db:"title"`

	// Description описание события.
	Description string `db:"description"`

	// StartAt дата и время начала события.
	StartAt time.Time `db:"start_at"`

	// EndAt дата и время окончания события.
	EndAt time.Time `db:"end_at"`

	// UserID идентификатор пользователя (владельца события).
	UserID uuid.UUID `db:"user_id"`

	// NotificationDuration за какое количество времени уведомить о начале события.
	NotificationDuration time.Duration `db:"notification_duration"`
}

// EventFilter предоставляет фильтр для поиска.
type EventFilter struct {
	// UserID идентификатор пользователя.
	UserID uuid.UUID

	// From дата и время начала события.
	From time.Time

	// To дата и время окончания события.
	To time.Time
}
