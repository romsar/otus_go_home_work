package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/proto/event"
)

type Server struct {
	r calendar.Repository
	event.UnimplementedEventServiceServer
}

func New(r calendar.Repository) *Server {
	return &Server{
		r: r,
	}
}

var _ event.EventServiceServer = (*Server)(nil)

func (s *Server) CreateEventV1(ctx context.Context, req *event.CreateEventRequestV1) (*event.EventResponseV1, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	e, err := s.r.CreateEvent(ctx, &calendar.Event{
		Title:                req.GetTitle(),
		Description:          req.GetDescription(),
		StartAt:              time.Unix(req.GetStartAt(), 0),
		EndAt:                time.Unix(req.GetEndAt(), 0),
		UserID:               userID,
		NotificationDuration: req.GetNotificationDuration(),
	})
	if err != nil {
		if errors.Is(err, calendar.ErrDateBusy) {
			return nil, status.Error(codes.InvalidArgument, "that date is already taken by another event")
		}

		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &event.EventResponseV1{
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

func (s *Server) UpdateEventV1(ctx context.Context, req *event.UpdateEventRequestV1) (*event.EventResponseV1, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	ID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid")
	}

	e, err := s.r.UpdateEvent(ctx, ID, &calendar.Event{
		Title:                req.GetTitle(),
		Description:          req.GetDescription(),
		StartAt:              time.Unix(req.GetStartAt(), 0),
		EndAt:                time.Unix(req.GetEndAt(), 0),
		UserID:               userID,
		NotificationDuration: req.GetNotificationDuration(),
	})
	if err != nil {
		if errors.Is(err, calendar.ErrDateBusy) {
			return nil, status.Error(codes.InvalidArgument, "that date is already taken by another event")
		}

		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &event.EventResponseV1{
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

func (s *Server) DeleteEventV1(ctx context.Context, req *event.DeleteEventRequestV1) (*emptypb.Empty, error) {
	ID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid")
	}

	if err := s.r.DeleteEvent(ctx, ID); err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) GetEventsForDayV1(ctx context.Context, req *event.GetEventsForDayRequestV1) (*event.EventsResponseV1, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	date, err := time.Parse("2006-01-02", req.GetDate())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date")
	}

	nextDay := date.AddDate(0, 0, 1)

	year, month, day := date.Date()
	nYear, nMonth, nDay := nextDay.Date()
	from := time.Date(year, month, day, 0, 0, 0, 0, date.UTC().Location())
	to := time.Date(nYear, nMonth, nDay, 0, 0, 0, 0, nextDay.UTC().Location())

	events, err := s.r.FindEvents(ctx, calendar.EventFilter{
		UserID: userID,
		From:   from,
		To:     to,
	})
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
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

	return &event.EventsResponseV1{
		Events: result,
	}, nil
}

func (s *Server) GetEventsForWeekV1(ctx context.Context, req *event.GetEventsForWeekRequestV1) (*event.EventsResponseV1, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	date, err := time.Parse("2006-01-02", req.GetStartDate())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date")
	}

	nextWeek := date.AddDate(0, 0, 7)

	year, month, day := date.Date()
	nYear, nMonth, nDay := nextWeek.Date()
	from := time.Date(year, month, day, 0, 0, 0, 0, date.UTC().Location())
	to := time.Date(nYear, nMonth, nDay, 0, 0, 0, 0, nextWeek.UTC().Location())

	events, err := s.r.FindEvents(ctx, calendar.EventFilter{
		UserID: userID,
		From:   from,
		To:     to,
	})
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
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

	return &event.EventsResponseV1{
		Events: result,
	}, nil
}

func (s *Server) GetEventsForMonthV1(ctx context.Context, req *event.GetEventsForMonthRequestV1) (*event.EventsResponseV1, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	date, err := time.Parse("2006-01-02", req.GetStartDate())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date")
	}

	nextMonth := date.AddDate(0, 1, 0)

	year, month, day := date.Date()
	nYear, nMonth, nDay := nextMonth.Date()
	from := time.Date(year, month, day, 0, 0, 0, 0, date.UTC().Location())
	to := time.Date(nYear, nMonth, nDay, 0, 0, 0, 0, nextMonth.UTC().Location())

	events, err := s.r.FindEvents(ctx, calendar.EventFilter{
		UserID: userID,
		From:   from,
		To:     to,
	})
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
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

	return &event.EventsResponseV1{
		Events: result,
	}, nil
}
