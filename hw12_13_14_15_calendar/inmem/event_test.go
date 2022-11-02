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
		startAt := mustParseDateTime("2022-12-05 10:00:00")
		endAt := mustParseDateTime("2022-12-05 10:00:30")
		notificationDuration := (10 * time.Hour).Minutes()

		event, err := repo.CreateEvent(ctx, &calendar.Event{
			ID:                   id,
			Title:                title,
			Description:          descr,
			StartAt:              startAt,
			EndAt:                endAt,
			UserID:               userID,
			NotificationDuration: uint32(notificationDuration),
		})

		require.NoError(t, err)
		require.NotNil(t, event)

		// check fields
		require.NotEqual(t, id, event.ID)
		require.Equal(t, title, event.Title)
		require.Equal(t, descr, event.Description)
		require.Equal(t, startAt, event.StartAt)
		require.Equal(t, endAt, event.EndAt)
		require.Equal(t, userID, event.UserID)
		require.Equal(t, uint32(notificationDuration), event.NotificationDuration)

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
		startAt := mustParseDateTime("2022-12-05 10:00:00")
		endAt := mustParseDateTime("2022-12-05 10:00:30")
		notificationDuration := (10 * time.Hour).Minutes()

		event, err = repo.UpdateEvent(ctx, id, &calendar.Event{
			Title:                title,
			Description:          descr,
			StartAt:              startAt,
			EndAt:                endAt,
			UserID:               userID,
			NotificationDuration: uint32(notificationDuration),
		})

		require.NoError(t, err)
		require.NotNil(t, event)

		// check fields
		require.Equal(t, id, event.ID)
		require.Equal(t, title, event.Title)
		require.Equal(t, descr, event.Description)
		require.Equal(t, startAt, event.StartAt)
		require.Equal(t, endAt, event.EndAt)
		require.Equal(t, userID, event.UserID)
		require.Equal(t, uint32(notificationDuration), event.NotificationDuration)

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

	t.Run("multiple", func(t *testing.T) {
		ctx := context.Background()
		repo := New()

		event1, err := repo.CreateEvent(ctx, &calendar.Event{})
		require.NoError(t, err)

		event2, err := repo.CreateEvent(ctx, &calendar.Event{})
		require.NoError(t, err)

		err = repo.DeleteEvent(ctx, event1.ID, event2.ID)
		require.NoError(t, err)

		// check storage
		require.Empty(t, repo.events)
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
		startAt := mustParseDateTime("2022-12-05 10:00:00")
		endAt := mustParseDateTime("2022-12-05 10:00:30")
		notificationDuration := (10 * time.Hour).Minutes()

		event, err := repo.CreateEvent(ctx, &calendar.Event{
			Title:                title,
			Description:          descr,
			StartAt:              startAt,
			EndAt:                endAt,
			UserID:               userID,
			NotificationDuration: uint32(notificationDuration),
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
		require.Equal(t, endAt, events[0].EndAt)
		require.Equal(t, userID, events[0].UserID)
		require.Equal(t, uint32(notificationDuration), events[0].NotificationDuration)
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
				name: "from match (equal)",
				args: args{
					event: calendar.Event{
						StartAt: mustParseDateTime("2022-05-10 15:20:30"),
					},
					filter: calendar.EventFilter{
						From: mustParseDateTime("2022-05-10 15:20:30"),
					},
				},
				found: true,
			},
			{
				name: "from match (greater)",
				args: args{
					event: calendar.Event{
						StartAt: mustParseDateTime("2022-05-10 15:20:40"),
					},
					filter: calendar.EventFilter{
						From: mustParseDateTime("2022-05-10 15:20:30"),
					},
				},
				found: true,
			},
			{
				name: "from skip",
				args: args{
					event: calendar.Event{
						StartAt: mustParseDateTime("2022-05-10 15:20:30"),
					},
					filter: calendar.EventFilter{
						From: mustParseDateTime("2022-05-10 15:20:31"),
					},
				},
				found: false,
			},

			// to filter
			{
				name: "to match (equal)",
				args: args{
					event: calendar.Event{
						StartAt: mustParseDateTime("2022-05-10 15:20:30"),
						EndAt:   mustParseDateTime("2022-05-10 15:50:30"),
					},
					filter: calendar.EventFilter{
						To: mustParseDateTime("2022-05-10 15:50:30"),
					},
				},
				found: true,
			},
			{
				name: "to match (less)",
				args: args{
					event: calendar.Event{
						StartAt: mustParseDateTime("2022-05-10 15:20:30"),
						EndAt:   mustParseDateTime("2022-05-10 15:50:20"),
					},
					filter: calendar.EventFilter{
						To: mustParseDateTime("2022-05-10 15:50:30"),
					},
				},
				found: true,
			},
			{
				name: "to skip",
				args: args{
					event: calendar.Event{
						StartAt: mustParseDateTime("2022-05-10 15:20:30"),
						EndAt:   mustParseDateTime("2022-05-10 15:50:30"),
					},
					filter: calendar.EventFilter{
						To: mustParseDateTime("2022-05-10 15:50:29"),
					},
				},
				found: false,
			},

			// is notified filter
			{
				name: "is notified match",
				args: args{
					event: calendar.Event{
						IsNotified: false,
					},
					filter: calendar.EventFilter{
						NotNotified: true,
					},
				},
				found: true,
			},
			{
				name: "is notified skip",
				args: args{
					event: calendar.Event{
						IsNotified: true,
					},
					filter: calendar.EventFilter{
						NotNotified: true,
					},
				},
				found: false,
			},

			// notify time
			{
				name: "notify time match eq",
				args: args{
					event: calendar.Event{
						StartAt:              mustParseDateTime("2022-05-10 16:00:00"),
						NotificationDuration: 30,
					},
					filter: calendar.EventFilter{
						NotifyTime: true,
					},
				},
				found: true,
			},
			{
				name: "notify time match gr",
				args: args{
					event: calendar.Event{
						StartAt:              mustParseDateTime("2022-05-10 15:59:59"),
						NotificationDuration: 30,
					},
					filter: calendar.EventFilter{
						NotifyTime: true,
					},
				},
				found: true,
			},
			{
				name: "notify time skip (already started)",
				args: args{
					event: calendar.Event{
						StartAt:              mustParseDateTime("2022-05-10 16:00:01"),
						NotificationDuration: 30,
					},
					filter: calendar.EventFilter{
						NotifyTime: true,
					},
				},
				found: false,
			},
			{
				name: "notify time skip (not yet)",
				args: args{
					event: calendar.Event{
						StartAt:              mustParseDateTime("2022-05-10 16:00:00"),
						NotificationDuration: 29,
					},
					filter: calendar.EventFilter{
						NotifyTime: true,
					},
				},
				found: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				timeNowFunc = func() time.Time {
					return mustParseDateTime("2022-05-10 15:30:00")
				}

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
		startAt := mustParseDateTime("2022-12-05 10:00:00")
		endAt := mustParseDateTime("2022-12-05 10:00:30")
		notificationDuration := (10 * time.Hour).Minutes()

		event, err := repo.CreateEvent(ctx, &calendar.Event{
			Title:                title,
			Description:          descr,
			StartAt:              startAt,
			EndAt:                endAt,
			UserID:               userID,
			NotificationDuration: uint32(notificationDuration),
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
		require.Equal(t, endAt, event.EndAt)
		require.Equal(t, userID, event.UserID)
		require.Equal(t, uint32(notificationDuration), event.NotificationDuration)
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

func TestRepository_checkDateBusy(t *testing.T) {
	tests := []struct {
		name       string
		events     []*calendar.Event
		ignoredIDs []uuid.UUID
		wantErr    error
	}{
		{
			name: "other event end when needle event begin",
			events: []*calendar.Event{
				{
					StartAt: mustParseDateTime("2022-05-10 14:50:00"),
					EndAt:   mustParseDateTime("2022-05-10 15:00:00"),
					UserID:  uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"),
				},
			},
			wantErr: nil,
		},
		{
			name: "other event end after needle event begin",
			events: []*calendar.Event{
				{
					StartAt: mustParseDateTime("2022-05-10 14:50:00"),
					EndAt:   mustParseDateTime("2022-05-10 15:00:01"),
					UserID:  uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"),
				},
			},
			wantErr: calendar.ErrDateBusy,
		},
		{
			name: "other event begin after needle event end",
			events: []*calendar.Event{
				{
					StartAt: mustParseDateTime("2022-05-10 15:30:00"),
					EndAt:   mustParseDateTime("2022-05-10 16:00:00"),
					UserID:  uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"),
				},
			},
			wantErr: nil,
		},
		{
			name: "other event begin before needle event end",
			events: []*calendar.Event{
				{
					StartAt: mustParseDateTime("2022-05-10 15:29:59"),
					EndAt:   mustParseDateTime("2022-05-10 15:59:59"),
					UserID:  uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"),
				},
			},
			wantErr: calendar.ErrDateBusy,
		},
		{
			name: "other event begin inside needle event",
			events: []*calendar.Event{
				{
					StartAt: mustParseDateTime("2022-05-10 15:20:00"),
					EndAt:   mustParseDateTime("2022-05-10 15:21:00"),
					UserID:  uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"),
				},
			},
			wantErr: calendar.ErrDateBusy,
		},
		{
			name: "ignore current event id",
			events: []*calendar.Event{
				{
					ID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					StartAt: mustParseDateTime("2022-05-10 15:20:00"),
					EndAt:   mustParseDateTime("2022-05-10 15:21:00"),
					UserID:  uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"),
				},
			},
			ignoredIDs: []uuid.UUID{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")},
			wantErr:    nil,
		},
		{
			name: "other user",
			events: []*calendar.Event{
				{
					StartAt: mustParseDateTime("2022-05-10 15:20:00"),
					EndAt:   mustParseDateTime("2022-05-10 15:21:00"),
					UserID:  uuid.New(),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := New()

			for _, e := range tt.events {
				if e.ID == uuid.Nil {
					e.ID = uuid.New()
				}
				repo.events[e.ID] = e
			}

			event := &calendar.Event{
				ID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				StartAt: mustParseDateTime("2022-05-10 15:00:00"),
				EndAt:   mustParseDateTime("2022-05-10 15:30:00"),
				UserID:  uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"),
			}

			err := repo.checkDateBusy(event, tt.ignoredIDs...)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func mustParseDateTime(str string) time.Time {
	dt, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		panic(err)
	}
	return dt
}
