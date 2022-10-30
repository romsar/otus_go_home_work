package inmem

import (
	"sync"

	"github.com/google/uuid"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
)

// Key обозначает ключ in-memory БД драйвера.
const Key = "inmemory"

// eventsMap определяет тип данных для in-memory хранилища событий.
type eventsMap map[uuid.UUID]*calendar.Event

// Repository реализует in-memory хранилище.
type Repository struct {
	eventMu sync.Mutex
	events  eventsMap
}

// New создает in-memory хранилище.
func New() *Repository {
	return &Repository{
		events: make(eventsMap),
	}
}
