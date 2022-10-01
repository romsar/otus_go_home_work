package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/mocks"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/proto/event"
)

func TestServer_CreateEventV1(t *testing.T) {
	t.Run("base test", func(t *testing.T) {
		m := mocks.NewModel(t)
		defer m.AssertExpectations(t)

		req := &event.CreateEventRequestV1{
			Event: &event.EventV1{
				Title:                "foo",
				Description:          "bar",
				StartAt:              int64(1664643702),
				EndAt:                int64(1664644150),
				UserId:               "123e4567-e89b-12d3-a456-426614174000",
				NotificationDuration: 30,
			},
		}

		m.On("CreateEvent", mock.Anything, &calendar.Event{
			Title:                "foo",
			Description:          "bar",
			StartAt:              time.Unix(1664643702, 0),
			EndAt:                time.Unix(1664644150, 0),
			UserID:               uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			NotificationDuration: 30,
		}).Return(&calendar.Event{
			ID:                   uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"),
			Title:                "foo",
			Description:          "bar",
			StartAt:              time.Unix(1664643702, 0),
			EndAt:                time.Unix(1664644150, 0),
			UserID:               uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			NotificationDuration: 30,
		}, nil).Once()

		s := Server{m: m}
		got, err := s.CreateEventV1(context.Background(), req)

		require.NoError(t, err)
		require.Equal(t, &event.EventReplyV1{
			Event: &event.EventV1{
				Id:                   "ef0d2079-e9a2-4810-8cae-eb6729c50580",
				Title:                "foo",
				Description:          "bar",
				StartAt:              1664643702,
				EndAt:                1664644150,
				UserId:               "123e4567-e89b-12d3-a456-426614174000",
				NotificationDuration: 30,
			},
		}, got)
	})
}

func TestServer_UpdateEventV1(t *testing.T) {
	t.Run("base test", func(t *testing.T) {
		m := mocks.NewModel(t)
		defer m.AssertExpectations(t)

		req := &event.UpdateEventRequestV1{
			Id: "ef0d2079-e9a2-4810-8cae-eb6729c50580",
			Event: &event.EventV1{
				Title:                "foo",
				Description:          "bar",
				StartAt:              int64(1664643702),
				EndAt:                int64(1664644150),
				UserId:               "123e4567-e89b-12d3-a456-426614174000",
				NotificationDuration: 30,
			},
		}

		m.On("UpdateEvent", mock.Anything, uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"), &calendar.Event{
			Title:                "foo",
			Description:          "bar",
			StartAt:              time.Unix(1664643702, 0),
			EndAt:                time.Unix(1664644150, 0),
			UserID:               uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			NotificationDuration: 30,
		}).Return(&calendar.Event{
			ID:                   uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"),
			Title:                "foo",
			Description:          "bar",
			StartAt:              time.Unix(1664643702, 0),
			EndAt:                time.Unix(1664644150, 0),
			UserID:               uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			NotificationDuration: 30,
		}, nil).Once()

		s := Server{m: m}
		got, err := s.UpdateEventV1(context.Background(), req)

		require.NoError(t, err)
		require.Equal(t, &event.EventReplyV1{
			Event: &event.EventV1{
				Id:                   "ef0d2079-e9a2-4810-8cae-eb6729c50580",
				Title:                "foo",
				Description:          "bar",
				StartAt:              1664643702,
				EndAt:                1664644150,
				UserId:               "123e4567-e89b-12d3-a456-426614174000",
				NotificationDuration: 30,
			},
		}, got)
	})
}

func TestServer_DeleteEventV1(t *testing.T) {
	t.Run("base test", func(t *testing.T) {
		m := mocks.NewModel(t)
		defer m.AssertExpectations(t)

		req := &event.DeleteEventRequestV1{
			Id: "ef0d2079-e9a2-4810-8cae-eb6729c50580",
		}

		m.On("DeleteEvent", mock.Anything, uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580")).
			Return(nil).
			Once()

		s := Server{m: m}
		got, err := s.DeleteEventV1(context.Background(), req)

		require.NoError(t, err)
		require.Equal(t, &event.DeleteEventReplyV1{
			Message: "Событие было успешно удалено!",
		}, got)
	})
}

func TestServer_GetEventsForDayV1(t *testing.T) {
	t.Run("base test", func(t *testing.T) {
		m := mocks.NewModel(t)
		defer m.AssertExpectations(t)

		req := &event.GetEventsForDayRequestV1{
			UserId: "123e4567-e89b-12d3-a456-426614174000",
			Date:   "2022-10-01",
		}

		m.On("FindEvents", mock.Anything, mock.Anything).Return([]*calendar.Event{
			{
				ID:                   uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"),
				Title:                "foo",
				Description:          "bar",
				StartAt:              time.Unix(1664643702, 0),
				EndAt:                time.Unix(1664644150, 0),
				UserID:               uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				NotificationDuration: 30,
			},
		}, nil).Once()

		s := Server{m: m}
		got, err := s.GetEventsForDayV1(context.Background(), req)

		require.NoError(t, err)
		require.Equal(t, &event.EventsReplyV1{
			Events: []*event.EventV1{
				{
					Id:                   "ef0d2079-e9a2-4810-8cae-eb6729c50580",
					Title:                "foo",
					Description:          "bar",
					StartAt:              1664643702,
					EndAt:                1664644150,
					UserId:               "123e4567-e89b-12d3-a456-426614174000",
					NotificationDuration: 30,
				},
			},
		}, got)
	})
}

func TestServer_GetEventsForWeekV1(t *testing.T) {
	t.Run("base test", func(t *testing.T) {
		m := mocks.NewModel(t)
		defer m.AssertExpectations(t)

		req := &event.GetEventsForWeekRequestV1{
			UserId:    "123e4567-e89b-12d3-a456-426614174000",
			StartDate: "2022-10-01",
		}

		m.On("FindEvents", mock.Anything, mock.Anything).Return([]*calendar.Event{
			{
				ID:                   uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"),
				Title:                "foo",
				Description:          "bar",
				StartAt:              time.Unix(1664643702, 0),
				EndAt:                time.Unix(1664644150, 0),
				UserID:               uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				NotificationDuration: 30,
			},
		}, nil).Once()

		s := Server{m: m}
		got, err := s.GetEventsForWeekV1(context.Background(), req)

		require.NoError(t, err)
		require.Equal(t, &event.EventsReplyV1{
			Events: []*event.EventV1{
				{
					Id:                   "ef0d2079-e9a2-4810-8cae-eb6729c50580",
					Title:                "foo",
					Description:          "bar",
					StartAt:              1664643702,
					EndAt:                1664644150,
					UserId:               "123e4567-e89b-12d3-a456-426614174000",
					NotificationDuration: 30,
				},
			},
		}, got)
	})
}

func TestServer_GetEventsForMonthV1(t *testing.T) {
	t.Run("base test", func(t *testing.T) {
		m := mocks.NewModel(t)
		defer m.AssertExpectations(t)

		req := &event.GetEventsForMonthRequestV1{
			UserId:    "123e4567-e89b-12d3-a456-426614174000",
			StartDate: "2022-10-01",
		}

		m.On("FindEvents", mock.Anything, mock.Anything).Return([]*calendar.Event{
			{
				ID:                   uuid.MustParse("ef0d2079-e9a2-4810-8cae-eb6729c50580"),
				Title:                "foo",
				Description:          "bar",
				StartAt:              time.Unix(1664643702, 0),
				EndAt:                time.Unix(1664644150, 0),
				UserID:               uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				NotificationDuration: 30,
			},
		}, nil).Once()

		s := Server{m: m}
		got, err := s.GetEventsForMonthV1(context.Background(), req)

		require.NoError(t, err)
		require.Equal(t, &event.EventsReplyV1{
			Events: []*event.EventV1{
				{
					Id:                   "ef0d2079-e9a2-4810-8cae-eb6729c50580",
					Title:                "foo",
					Description:          "bar",
					StartAt:              1664643702,
					EndAt:                1664644150,
					UserId:               "123e4567-e89b-12d3-a456-426614174000",
					NotificationDuration: 30,
				},
			},
		}, got)
	})
}
