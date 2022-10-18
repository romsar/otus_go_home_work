package postgres

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
)

// CreateEvent создать событие.
func (repo *Repository) CreateEvent(ctx context.Context, e *calendar.Event) (*calendar.Event, error) {
	if err := repo.checkDateBusy(ctx, e); err != nil {
		return nil, errors.Wrap(err, "create event")
	}

	e.ID = uuid.New()

	event := new(calendar.Event)
	err := repo.db.QueryRowxContext(
		ctx,
		`INSERT INTO events (id, title, description, start_at, end_at, user_id, notification_duration) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;`, //nolint:lll
		e.ID, e.Title, e.Description, e.StartAt, e.EndAt, e.UserID, e.NotificationDuration,
	).StructScan(event)
	if err != nil {
		return nil, errors.Wrap(err, "create event")
	}

	return event, nil
}

// UpdateEvent обновить событие.
func (repo *Repository) UpdateEvent(ctx context.Context, id uuid.UUID, e *calendar.Event) (*calendar.Event, error) {
	if _, err := repo.findEventByID(ctx, id); err != nil {
		return nil, errors.Wrap(err, "update event")
	}

	if err := repo.checkDateBusy(ctx, e); err != nil {
		return nil, errors.Wrap(err, "update event")
	}

	event := new(calendar.Event)
	err := repo.db.QueryRowxContext(
		ctx,
		`UPDATE events SET title = $1, description = $2, start_at = $3, end_at = $4, user_id = $5, notification_duration = $6 WHERE id = $7 RETURNING *;`, //nolint:lll
		e.Title, e.Description, e.StartAt, e.EndAt, e.UserID, e.NotificationDuration, id,
	).StructScan(event)
	if err != nil {
		return nil, errors.Wrap(err, "update event")
	}

	return event, nil
}

// DeleteEvent удалить событие.
func (repo *Repository) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	if _, err := repo.findEventByID(ctx, id); err != nil {
		return errors.Wrap(err, "delete event")
	}

	_, err := repo.db.ExecContext(ctx, `DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		return errors.Wrap(err, "delete event")
	}

	return nil
}

// FindEvents найти множество событий.
func (repo *Repository) FindEvents(ctx context.Context, filter calendar.EventFilter) ([]*calendar.Event, error) {
	where, args, counter := []string{"1 = 1"}, []interface{}{}, 1

	if filter.UserID != uuid.Nil {
		where, args = append(where, "user_id = $"+strconv.Itoa(counter)), append(args, filter.UserID)
		counter++
	}

	if !filter.From.IsZero() {
		where, args = append(where, "start_at >= $"+strconv.Itoa(counter)), append(args, filter.From)
		counter++
	}

	if !filter.To.IsZero() {
		where, args = append(where, `end_at <= $`+strconv.Itoa(counter)), append(args, filter.To)
		counter++ //nolint:ineffassign,wastedassign
	}

	events := make([]*calendar.Event, 0)

	err := repo.db.SelectContext(ctx, &events, `
		SELECT * FROM events
		WHERE `+strings.Join(where, " AND "),
		args...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "find events")
	}

	return events, nil
}

// FindEventByID найти событие по его идентификатору.
func (repo *Repository) FindEventByID(ctx context.Context, id uuid.UUID) (*calendar.Event, error) {
	event, err := repo.findEventByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "find event")
	}

	return event, nil
}

// findEventByID найти событие по его идентификатору.
func (repo *Repository) findEventByID(ctx context.Context, id uuid.UUID) (*calendar.Event, error) {
	event := new(calendar.Event)

	err := repo.db.GetContext(ctx, event, `SELECT * FROM events WHERE id=$1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = calendar.ErrNotFound
		}

		return nil, err
	}

	return event, nil
}

// checkDateBusy проверка на свободное время.
// Если время занято, то вернет ошибку calendar.ErrDateBusy.
func (repo *Repository) checkDateBusy(ctx context.Context, event *calendar.Event) error {
	query := `
			SELECT count(*) AS count
			FROM events
			WHERE user_id = $1
			  AND id != $2
			  AND ((start_at >= $3 AND start_at < $4) OR (end_at > $3 AND end_at <= $4))
		`

	var count int
	err := repo.db.QueryRowContext(
		ctx, query, event.UserID, event.ID,
		event.StartAt, event.EndAt,
	).Scan(&count)
	if err != nil {
		return errors.Wrap(err, "check date busy")
	}

	if count > 0 {
		return calendar.ErrDateBusy
	}

	return nil
}
