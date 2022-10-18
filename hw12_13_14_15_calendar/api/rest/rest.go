package rest

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
)

// Config декларирует настройки REST.
type Config struct {
	// Address адрес, по которому будет доступен сервер.
	Address string
}

// Server декларирует REST сервер.
type Server struct {
	httpServer *http.Server
	m          Model
}

// Model декларирует контракт модели.
type Model interface {
	// CreateEvent создать событие.
	CreateEvent(ctx context.Context, e *calendar.Event) (*calendar.Event, error)

	// UpdateEvent обновить событие.
	UpdateEvent(ctx context.Context, id uuid.UUID, e *calendar.Event) (*calendar.Event, error)

	// DeleteEvent удалить событие.
	DeleteEvent(ctx context.Context, id uuid.UUID) error

	// FindEvents найти множество событий.
	FindEvents(ctx context.Context, filter calendar.EventFilter) ([]*calendar.Event, error)

	// FindEventByID найти событие по его идентификатору.
	FindEventByID(ctx context.Context, id uuid.UUID) (*calendar.Event, error)
}

// New инициализирует REST.
func New(cfg Config, m Model) Server {
	return Server{
		httpServer: &http.Server{
			Addr: cfg.Address,
		},
		m: m,
	}
}

// Listen запускает сервер.
func (s Server) Listen(handlers ...func(next http.Handler) http.Handler) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.hello)

	s.httpServer.Handler = mux

	for _, h := range handlers {
		s.httpServer.Handler = h(s.httpServer.Handler)
	}

	if err := s.httpServer.ListenAndServe(); err != nil {
		return errors.Wrap(err, "server listen")
	}

	return nil
}

// Close завершает работу сервера.
func (s Server) Close(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s Server) hello(w http.ResponseWriter, r *http.Request) {
	// todo s.m.SomeMethod()

	w.Write([]byte("Hello, world!"))
}
