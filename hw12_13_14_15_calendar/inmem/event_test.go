package inmem

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
)

//nolint:goconst
func TestEventStorage_CreateEvent(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		repo := New()

		id := uuid.New()
		title := "foo"
		descr := "bar"
		userID := uuid.New()
		startAt, _ := time.Parse("2006-01-02", "2022-12-05")
		duration := 30 * time.Second
		notificationDuration := 10 * time.Hour

		event, err := repo.CreateEvent(ctx, &calendar.Event{
			ID:                   id,
			Title:                title,
			Description:          descr,
			StartAt:              startAt,
			Duration:             duration,
			UserID:               userID,
			NotificationDuration: notificationDuration,
		})

		require.NoError(t, err)
		require.NotNil(t, event)

		// check fields
		require.NotEqual(t, id, event.ID)
		require.Equal(t, title, event.Title)
		require.Equal(t, descr, event.Description)
		require.Equal(t, startAt, event.StartAt)
		require.Equal(t, duration, event.Duration)
		require.Equal(t, userID, event.UserID)
		require.Equal(t, notificationDuration, event.NotificationDuration)

		// check storage
		require.Len(t, repo.events, 1)

		event, err = repo.findEventByID(event.ID)
		require.NotNil(t, event)
		require.NoError(t, err)

		event, err = repo.findEventByID(id)
		require.Nil(t, event)
		require.ErrorIs(t, err, calendar.ErrNotFound)
	})
}

func TestEventStorage_UpdateEvent(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		repo := New()

		event, err := repo.CreateEvent(ctx, &calendar.Event{})
		require.NoError(t, err)

		id := event.ID
		title := "foo"
		descr := "bar"
		userID := uuid.New()
		startAt, _ := time.Parse("2006-01-02", "2022-12-05")
		duration := 30 * time.Second
		notificationDuration := 10 * time.Hour

		event, err = repo.UpdateEvent(ctx, id, &calendar.Event{
			Title:                title,
			Description:          descr,
			StartAt:              startAt,
			Duration:             duration,
			UserID:               userID,
			NotificationDuration: notificationDuration,
		})

		require.NoError(t, err)
		require.NotNil(t, event)

		// check fields
		require.Equal(t, id, event.ID)
		require.Equal(t, title, event.Title)
		require.Equal(t, descr, event.Description)
		require.Equal(t, startAt, event.StartAt)
		require.Equal(t, duration, event.Duration)
		require.Equal(t, userID, event.UserID)
		require.Equal(t, notificationDuration, event.NotificationDuration)

		// check storage
		require.Len(t, repo.events, 1)
	})

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()
		repo := New()

		id := uuid.New()

		event, err := repo.UpdateEvent(ctx, id, &calendar.Event{})

		require.ErrorIs(t, err, calendar.ErrNotFound)
		require.Nil(t, event)

		// check storage
		require.Empty(t, repo.events)
	})
}

func TestEventStorage_DeleteEvent(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		repo := New()

		event, err := repo.CreateEvent(ctx, &calendar.Event{})
		require.NoError(t, err)

		id := event.ID

		err = repo.DeleteEvent(ctx, id)

		require.NoError(t, err)

		// check storage
		require.Empty(t, repo.events)
	})

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()
		repo := New()

		id := uuid.New()

		err := repo.DeleteEvent(ctx, id)

		require.ErrorIs(t, err, calendar.ErrNotFound)
	})
}

//nolint:funlen
func TestEventStorage_FindEvents(t *testing.T) {
	t.Parallel()

	t.Run("format", func(t *testing.T) {
		ctx := context.Background()
		repo := New()

		title := "foo"
		descr := "bar"
		userID := uuid.New()
		startAt, _ := time.Parse("2006-01-02", "2022-12-05")
		duration := 30 * time.Second
		notificationDuration := 10 * time.Hour

		event, err := repo.CreateEvent(ctx, &calendar.Event{
			Title:                title,
			Description:          descr,
			StartAt:              startAt,
			Duration:             duration,
			UserID:               userID,
			NotificationDuration: notificationDuration,
		})
		require.NoError(t, err)

		events, err := repo.FindEvents(ctx, calendar.EventFilter{})

		require.NoError(t, err)
		require.Len(t, events, 1)

		// check fields
		require.Equal(t, event.ID, events[0].ID)
		require.Equal(t, title, events[0].Title)
		require.Equal(t, descr, events[0].Description)
		require.Equal(t, startAt, events[0].StartAt)
		require.Equal(t, duration, events[0].Duration)
		require.Equal(t, userID, events[0].UserID)
		require.Equal(t, notificationDuration, events[0].NotificationDuration)
	})

	t.Run("filter", func(t *testing.T) {
		type args struct {
			event  calendar.Event
			filter calendar.EventFilter
		}
		tests := []struct {
			name  string
			args  args
			found bool
		}{
			// user id filter
			{
				name: "user id match",
				args: args{
					event: calendar.Event{
						UserID: uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"),
					},
					filter: calendar.EventFilter{
						UserID: uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"),
					},
				},
				found: true,
			},
			{
				name: "user id skip",
				args: args{
					event: calendar.Event{
						UserID: uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"),
					},
					filter: calendar.EventFilter{
						UserID: uuid.MustParse("2f0d2079-e9a2-4810-8cae-eb6729c50580"),
					},
				},
				found: false,
			},
			// from filter
			{
				name: "from match",
				args: args{
					event: calendar.Event{
						StartAt: func() time.Time {
							date, err := time.Parse("2006-01-02 15:04:05", "2022-05-10 15:20:30")
							require.NoError(t, err)
							return date
						}(),
					},
					filter: calendar.EventFilter{
						From: func() time.Time {
							date, err := time.Parse("2006-01-02 15:04:05", "2022-05-10 15:20:30")
							require.NoError(t, err)
							return date
						}(),
					},
				},
				found: true,
			},
			{
				name: "from skip",
				args: args{
					event: calendar.Event{
						StartAt: func() time.Time {
							date, err := time.Parse("2006-01-02 15:04:05", "2022-05-10 15:20:30")
							require.NoError(t, err)
							return date
						}(),
					},
					filter: calendar.EventFilter{
						From: func() time.Time {
							date, err := time.Parse("2006-01-02 15:04:05", "2022-05-10 15:20:29")
							require.NoError(t, err)
							return date
						}(),
					},
				},
				found: false,
			},

			// to filter
			{
				name: "to match",
				args: args{
					event: calendar.Event{
						StartAt: func() time.Time {
							date, err := time.Parse("2006-01-02 15:04:05", "2022-05-10 15:20:30")
							require.NoError(t, err)
							return date
						}(),
						Duration: 30 * time.Minute,
					},
					filter: calendar.EventFilter{
						To: func() time.Time {
							date, err := time.Parse("2006-01-02 15:04:05", "2022-05-10 15:50:30")
							require.NoError(t, err)
							return date
						}(),
					},
				},
				found: true,
			},
			{
				name: "to skip",
				args: args{
					event: calendar.Event{
						StartAt: func() time.Time {
							date, err := time.Parse("2006-01-02 15:04:05", "2022-05-10 15:20:30")
							require.NoError(t, err)
							return date
						}(),
						Duration: 30 * time.Minute,
					},
					filter: calendar.EventFilter{
						To: func() time.Time {
							date, err := time.Parse("2006-01-02 15:04:05", "2022-05-10 15:50:31")
							require.NoError(t, err)
							return date
						}(),
					},
				},
				found: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctx := context.Background()
				repo := New()

				_, err := repo.CreateEvent(ctx, &tt.args.event)
				require.NoError(t, err)

				events, err := repo.FindEvents(ctx, tt.args.filter)
				require.NoError(t, err)

				if tt.found {
					require.Len(t, events, 1)
				} else {
					require.Empty(t, events)
				}
			})
		}
	})
}

func TestEventStorage_FindEvent(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		repo := New()

		title := "foo"
		descr := "bar"
		userID := uuid.New()
		startAt, _ := time.Parse("2006-01-02", "2022-12-05")
		duration := 30 * time.Second
		notificationDuration := 10 * time.Hour

		event, err := repo.CreateEvent(ctx, &calendar.Event{
			Title:                title,
			Description:          descr,
			StartAt:              startAt,
			Duration:             duration,
			UserID:               userID,
			NotificationDuration: notificationDuration,
		})
		require.NoError(t, err)
		id := event.ID

		event, err = repo.FindEventByID(ctx, id)

		require.NoError(t, err)

		// check fields
		require.Equal(t, id, id)
		require.Equal(t, title, event.Title)
		require.Equal(t, descr, event.Description)
		require.Equal(t, startAt, event.StartAt)
		require.Equal(t, duration, event.Duration)
		require.Equal(t, userID, event.UserID)
		require.Equal(t, notificationDuration, event.NotificationDuration)
	})

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()
		repo := New()

		id := uuid.New()

		event, err := repo.FindEventByID(ctx, id)

		require.Nil(t, event)
		require.ErrorIs(t, err, calendar.ErrNotFound)
	})
}
