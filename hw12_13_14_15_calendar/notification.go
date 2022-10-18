package calendar

import (
	"time"

	"github.com/google/uuid"
)

// Notification (уведомление) - временная сущность,
// в БД не хранится, складывается в очередь для рассыльщика.
type Notification struct {
	// EventID идентификатор события.
	EventID uuid.UUID

	// EventTitle заголовок события.
	EventTitle string

	// EventStartAt дата и время начала события.
	EventStartAt time.Time

	// UserID пользователь, кому отправить уведомление.
	UserID uuid.UUID
}
