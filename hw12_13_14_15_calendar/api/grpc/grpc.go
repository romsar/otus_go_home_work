package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/proto/event"
)

type Server struct {
	m calendar.Model
	event.UnimplementedEventServiceServer
}

func New(m calendar.Model) Server {
	return Server{
		m: m,
	}
}

var _ event.EventServiceServer = (*Server)(nil)

func (s Server) CreateEventV1(ctx context.Context, req *event.CreateEventRequestV1) (*event.EventReplyV1, error) {
	userID, err := uuid.Parse(req.Event.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	e, err := s.m.CreateEvent(ctx, &calendar.Event{
		Title:                req.Event.Title,
		Description:          req.Event.Description,
		StartAt:              time.Unix(req.Event.StartAt, 0),
		EndAt:                time.Unix(req.Event.EndAt, 0),
		UserID:               userID,
		NotificationDuration: req.Event.NotificationDuration,
	})
	if err != nil {
		if errors.Is(err, calendar.ErrDateBusy) {
			return nil, status.Error(codes.InvalidArgument, "that date is already taken by another event")
		}

		return nil, status.Error(codes.Unavailable, "error while creating event")
	}

	return &event.EventReplyV1{
		Event: &event.EventV1{
			Id:                   e.ID.String(),
			Title:                e.Title,
			Description:          e.Description,
			StartAt:              e.StartAt.Unix(),
			EndAt:                e.EndAt.Unix(),
			UserId:               e.UserID.String(),
			NotificationDuration: e.NotificationDuration,
		},
	}, nil
}

func (s Server) UpdateEventV1(ctx context.Context, req *event.UpdateEventRequestV1) (*event.EventReplyV1, error) {
	userID, err := uuid.Parse(req.Event.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	ID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid")
	}

	e, err := s.m.UpdateEvent(ctx, ID, &calendar.Event{
		Title:                req.Event.Title,
		Description:          req.Event.Description,
		StartAt:              time.Unix(req.Event.StartAt, 0),
		EndAt:                time.Unix(req.Event.EndAt, 0),
		UserID:               userID,
		NotificationDuration: req.Event.NotificationDuration,
	})
	if err != nil {
		if errors.Is(err, calendar.ErrDateBusy) {
			return nil, status.Error(codes.InvalidArgument, "that date is already taken by another event")
		}

		return nil, status.Error(codes.Unavailable, "error while updating event")
	}

	return &event.EventReplyV1{
		Event: &event.EventV1{
			Id:                   e.ID.String(),
			Title:                e.Title,
			Description:          e.Description,
			StartAt:              e.StartAt.Unix(),
			EndAt:                e.EndAt.Unix(),
			UserId:               e.UserID.String(),
			NotificationDuration: e.NotificationDuration,
		},
	}, nil
}

func (s Server) DeleteEventV1(ctx context.Context, req *event.DeleteEventRequestV1) (*event.DeleteEventReplyV1, error) {
	ID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid")
	}

	if err := s.m.DeleteEvent(ctx, ID); err != nil {
		return nil, status.Error(codes.Unavailable, "error while deleting event")
	}

	return &event.DeleteEventReplyV1{
		Message: "Событие было успешно удалено!",
	}, nil
}

func (s Server) GetEventsForDayV1(ctx context.Context, req *event.GetEventsForDayRequestV1) (*event.EventsReplyV1, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date")
	}

	nextDay := date.AddDate(0, 0, 1)

	year, month, day := date.Date()
	nYear, nMonth, nDay := nextDay.Date()
	from := time.Date(year, month, day, 0, 0, 0, 0, date.UTC().Location())
	to := time.Date(nYear, nMonth, nDay, 0, 0, 0, 0, nextDay.UTC().Location())

	events, err := s.m.FindEvents(ctx, calendar.EventFilter{
		UserID: userID,
		From:   from,
		To:     to,
	})
	if err != nil {
		return nil, status.Error(codes.Unavailable, "error while getting events for day")
	}

	result := make([]*event.EventV1, 0, len(events))

	for _, e := range events {
		result = append(result, &event.EventV1{
			Id:                   e.ID.String(),
			Title:                e.Title,
			Description:          e.Description,
			StartAt:              e.StartAt.Unix(),
			EndAt:                e.EndAt.Unix(),
			UserId:               e.UserID.String(),
			NotificationDuration: e.NotificationDuration,
		})
	}

	return &event.EventsReplyV1{
		Events: result,
	}, nil
}

func (s Server) GetEventsForWeekV1(ctx context.Context, req *event.GetEventsForWeekRequestV1) (*event.EventsReplyV1, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	date, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date")
	}

	nextWeek := date.AddDate(0, 0, 7)

	year, month, day := date.Date()
	nYear, nMonth, nDay := nextWeek.Date()
	from := time.Date(year, month, day, 0, 0, 0, 0, date.UTC().Location())
	to := time.Date(nYear, nMonth, nDay, 0, 0, 0, 0, nextWeek.UTC().Location())

	events, err := s.m.FindEvents(ctx, calendar.EventFilter{
		UserID: userID,
		From:   from,
		To:     to,
	})
	if err != nil {
		return nil, status.Error(codes.Unavailable, "error while getting events for week")
	}

	result := make([]*event.EventV1, 0, len(events))

	for _, e := range events {
		result = append(result, &event.EventV1{
			Id:                   e.ID.String(),
			Title:                e.Title,
			Description:          e.Description,
			StartAt:              e.StartAt.Unix(),
			EndAt:                e.EndAt.Unix(),
			UserId:               e.UserID.String(),
			NotificationDuration: e.NotificationDuration,
		})
	}

	return &event.EventsReplyV1{
		Events: result,
	}, nil
}

func (s Server) GetEventsForMonthV1(ctx context.Context, req *event.GetEventsForMonthRequestV1) (*event.EventsReplyV1, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	date, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date")
	}

	nextMonth := date.AddDate(0, 1, 0)

	year, month, day := date.Date()
	nYear, nMonth, nDay := nextMonth.Date()
	from := time.Date(year, month, day, 0, 0, 0, 0, date.UTC().Location())
	to := time.Date(nYear, nMonth, nDay, 0, 0, 0, 0, nextMonth.UTC().Location())

	events, err := s.m.FindEvents(ctx, calendar.EventFilter{
		UserID: userID,
		From:   from,
		To:     to,
	})
	if err != nil {
		return nil, status.Error(codes.Unavailable, "error while getting events for month")
	}

	result := make([]*event.EventV1, 0, len(events))

	for _, e := range events {
		result = append(result, &event.EventV1{
			Id:                   e.ID.String(),
			Title:                e.Title,
			Description:          e.Description,
			StartAt:              e.StartAt.Unix(),
			EndAt:                e.EndAt.Unix(),
			UserId:               e.UserID.String(),
			NotificationDuration: e.NotificationDuration,
		})
	}

	return &event.EventsReplyV1{
		Events: result,
	}, nil
}
