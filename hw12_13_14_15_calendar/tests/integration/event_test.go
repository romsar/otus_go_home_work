//go:build integration
// +build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	_ "github.com/lib/pq"

	"github.com/RomanSarvarov/otus_go_home_work/calendar/proto/event"
)

type EventSuite struct {
	suite.Suite
	ctx context.Context

	grpcConn *grpc.ClientConn
	pgConn   *sqlx.DB

	eventClient event.EventServiceClient
}

func (s *EventSuite) SetupSuite() {
	// GRPC setup
	s.setupGRPC()

	// Postgres setup
	s.setupPG()

	s.ctx = context.Background()
}

func (s *EventSuite) SetupTest() {
	s.cleanTables("events")
}

func (s *EventSuite) TearDownSuite() {
	s.cleanTables("events")
	s.Require().NoError(s.grpcConn.Close())
	s.Require().NoError(s.pgConn.Close())
}

func (s *EventSuite) TestCreateEvent() {
	s.Run("success", func() {
		s.SetupTest()

		// insert event for different user
		_, err := s.pgConn.QueryContext(
			s.ctx,
			`INSERT INTO events (id, title, description, start_at, end_at, user_id, notification_duration, is_notified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`, //nolint:lll
			uuid.New(), "aaa", "bbb", time.Unix(1664643900, 0).UTC(), time.Unix(1664644000, 0).UTC(), "ef0d2079-29a2-4810-8cae-eb6729c50580", 30, false, //nolint:lll
		)
		s.Require().NoError(err)

		resp, err := s.eventClient.CreateEventV1(s.ctx, &event.CreateEventRequestV1{
			Event: &event.EventV1{
				Title:                "foo",
				Description:          "bar",
				StartAt:              1664643702,
				EndAt:                1664644150,
				UserId:               "ef0d2079-e9a2-4810-8cae-eb6729c50580",
				NotificationDuration: 25,
			},
		})
		s.Require().NoError(err)
		s.Require().NotEmpty(resp.Event.Id)
		s.Require().Equal("foo", resp.Event.Title)
		s.Require().Equal("bar", resp.Event.Description)
		s.Require().Equal(int64(1664643702), resp.Event.StartAt)
		s.Require().Equal(int64(1664644150), resp.Event.EndAt)
		s.Require().Equal("ef0d2079-e9a2-4810-8cae-eb6729c50580", resp.Event.UserId)
		s.Require().Equal(uint32(25), resp.Event.NotificationDuration)

		query := `
			SELECT count(*) AS count
			FROM events
			WHERE id = $1
		`

		var count int
		err = s.pgConn.QueryRowContext(s.ctx, query, resp.Event.Id).Scan(&count)
		s.Require().NoError(err)

		s.Require().Equal(1, count)
	})

	s.Run("busy", func() {
		s.SetupTest()

		_, err := s.pgConn.QueryContext(
			s.ctx,
			`INSERT INTO events (id, title, description, start_at, end_at, user_id, notification_duration, is_notified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`, //nolint:lll
			uuid.New(), "aaa", "bbb", time.Unix(1664643900, 0).UTC(), time.Unix(1664644000, 0).UTC(), "ef0d2079-e9a2-4810-8cae-eb6729c50580", 30, false, //nolint:lll
		)
		s.Require().NoError(err)

		_, err = s.eventClient.CreateEventV1(s.ctx, &event.CreateEventRequestV1{
			Event: &event.EventV1{
				Title:                "foo",
				Description:          "bar",
				StartAt:              1664643900,
				EndAt:                1664644000,
				UserId:               "ef0d2079-e9a2-4810-8cae-eb6729c50580",
				NotificationDuration: 25,
			},
		})
		s.Require().Equal(codes.InvalidArgument, status.Code(err))

		query := `
			SELECT count(*) AS count
			FROM events
		`

		var count int
		err = s.pgConn.QueryRowContext(s.ctx, query).Scan(&count)
		s.Require().NoError(err)

		s.Require().Equal(1, count)
	})

	s.Run("no user passed", func() {
		s.SetupTest()

		_, err := s.eventClient.CreateEventV1(s.ctx, &event.CreateEventRequestV1{
			Event: &event.EventV1{
				Title:                "foo",
				Description:          "bar",
				StartAt:              1664643900,
				EndAt:                1664644000,
				UserId:               "",
				NotificationDuration: 25,
			},
		})
		s.Require().Equal(codes.InvalidArgument, status.Code(err))

		query := `
			SELECT count(*) AS count
			FROM events
		`

		var count int
		err = s.pgConn.QueryRowContext(s.ctx, query).Scan(&count)
		s.Require().NoError(err)

		s.Require().Equal(0, count)
	})
}

func (s *EventSuite) TestGetEventsForDay() {
	s.Run("empty", func() {
		s.SetupTest()

		resp, err := s.eventClient.GetEventsForDayV1(s.ctx, &event.GetEventsForDayRequestV1{
			UserId: "ef0d2079-e9a2-4810-8cae-eb6729c50580",
			Date:   "2022-10-12",
		})

		s.Require().NoError(err)
		s.Require().Empty(resp.Events)
	})

	s.Run("different user", func() {
		s.SetupTest()

		_, err := s.pgConn.QueryContext(
			s.ctx,
			`INSERT INTO events (id, title, description, start_at, end_at, user_id, notification_duration, is_notified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`, //nolint:lll
			uuid.New(), "aaa", "bbb", time.Date(2022, 10, 12, 12, 30, 0, 0, time.UTC), time.Date(2022, 10, 12, 14, 30, 0, 0, time.UTC), "2f0d2079-e9a2-4810-8cae-eb6729c50580", 30, false, //nolint:lll
		)
		s.Require().NoError(err)

		resp, err := s.eventClient.GetEventsForDayV1(s.ctx, &event.GetEventsForDayRequestV1{
			UserId: "ef0d2079-e9a2-4810-8cae-eb6729c50580",
			Date:   "2022-10-12",
		})

		s.Require().NoError(err)
		s.Require().Empty(resp.Events)
	})

	s.Run("different day", func() {
		s.SetupTest()

		_, err := s.pgConn.QueryContext(
			s.ctx,
			`INSERT INTO events (id, title, description, start_at, end_at, user_id, notification_duration, is_notified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`, //nolint:lll
			uuid.New(), "aaa", "bbb", time.Date(2022, 10, 12, 12, 30, 0, 0, time.UTC), time.Date(2022, 10, 12, 14, 30, 0, 0, time.UTC), "ef0d2079-e9a2-4810-8cae-eb6729c50580", 30, false, //nolint:lll
		)
		s.Require().NoError(err)

		resp, err := s.eventClient.GetEventsForDayV1(s.ctx, &event.GetEventsForDayRequestV1{
			UserId: "ef0d2079-e9a2-4810-8cae-eb6729c50580",
			Date:   "2022-10-11",
		})

		s.Require().NoError(err)
		s.Require().Empty(resp.Events)
	})

	s.Run("success", func() {
		s.SetupTest()

		_, err := s.pgConn.QueryContext(
			s.ctx,
			`INSERT INTO events (id, title, description, start_at, end_at, user_id, notification_duration, is_notified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`, //nolint:lll
			uuid.New(), "aaa", "bbb", time.Date(2022, 10, 12, 12, 30, 0, 0, time.UTC), time.Date(2022, 10, 12, 14, 30, 0, 0, time.UTC), "ef0d2079-e9a2-4810-8cae-eb6729c50580", 30, false, //nolint:lll
		)
		s.Require().NoError(err)

		resp, err := s.eventClient.GetEventsForDayV1(s.ctx, &event.GetEventsForDayRequestV1{
			UserId: "ef0d2079-e9a2-4810-8cae-eb6729c50580",
			Date:   "2022-10-12",
		})

		s.Require().NoError(err)
		s.Require().Len(resp.Events, 1)

		e := resp.Events[0]

		s.Require().NotEmpty(e.Id)
		s.Require().Equal("aaa", e.Title)
		s.Require().Equal("bbb", e.Description)
		s.Require().Equal(int64(1665577800), e.StartAt)
		s.Require().Equal(int64(1665585000), e.EndAt)
		s.Require().Equal("ef0d2079-e9a2-4810-8cae-eb6729c50580", e.UserId)
		s.Require().Equal(uint32(30), e.NotificationDuration)
	})
}

func (s *EventSuite) TestGetEventsForWeek() {
	s.SetupTest()

	_, err := s.pgConn.QueryContext(
		s.ctx,
		`INSERT INTO events (id, title, description, start_at, end_at, user_id, notification_duration, is_notified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`, //nolint:lll
		uuid.New(), "aaa", "bbb", time.Date(2022, 10, 12, 12, 30, 0, 0, time.UTC), time.Date(2022, 10, 12, 14, 30, 0, 0, time.UTC), "ef0d2079-e9a2-4810-8cae-eb6729c50580", 30, false, //nolint:lll
	)
	s.Require().NoError(err)

	resp, err := s.eventClient.GetEventsForWeekV1(s.ctx, &event.GetEventsForWeekRequestV1{
		UserId:    "ef0d2079-e9a2-4810-8cae-eb6729c50580",
		StartDate: "2022-10-09",
	})

	s.Require().NoError(err)
	s.Require().Len(resp.Events, 1)

	e := resp.Events[0]

	s.Require().NotEmpty(e.Id)
	s.Require().Equal("aaa", e.Title)
	s.Require().Equal("bbb", e.Description)
	s.Require().Equal(int64(1665577800), e.StartAt)
	s.Require().Equal(int64(1665585000), e.EndAt)
	s.Require().Equal("ef0d2079-e9a2-4810-8cae-eb6729c50580", e.UserId)
	s.Require().Equal(uint32(30), e.NotificationDuration)
}

func (s *EventSuite) TestGetEventsForMonth() {
	_, err := s.pgConn.QueryContext(
		s.ctx,
		`INSERT INTO events (id, title, description, start_at, end_at, user_id, notification_duration, is_notified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`, //nolint:lll
		uuid.New(), "aaa", "bbb", time.Date(2022, 10, 12, 12, 30, 0, 0, time.UTC), time.Date(2022, 10, 12, 14, 30, 0, 0, time.UTC), "ef0d2079-e9a2-4810-8cae-eb6729c50580", 30, false, //nolint:lll
	)
	s.Require().NoError(err)

	resp, err := s.eventClient.GetEventsForMonthV1(s.ctx, &event.GetEventsForMonthRequestV1{
		UserId:    "ef0d2079-e9a2-4810-8cae-eb6729c50580",
		StartDate: "2022-09-20",
	})

	s.Require().NoError(err)
	s.Require().Len(resp.Events, 1)

	e := resp.Events[0]

	s.Require().NotEmpty(e.Id)
	s.Require().Equal("aaa", e.Title)
	s.Require().Equal("bbb", e.Description)
	s.Require().Equal(int64(1665577800), e.StartAt)
	s.Require().Equal(int64(1665585000), e.EndAt)
	s.Require().Equal("ef0d2079-e9a2-4810-8cae-eb6729c50580", e.UserId)
	s.Require().Equal(uint32(30), e.NotificationDuration)
}

func (s *EventSuite) TestSendEventNotification() {
	startAt := time.Now().Add(time.Hour * 24).UTC()
	endAt := startAt.Add(30 * time.Minute).UTC()

	_, err := s.pgConn.QueryContext(
		s.ctx,
		`INSERT INTO events (id, title, description, start_at, end_at, user_id, notification_duration, is_notified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`, //nolint:lll
		uuid.New(), "aaa", "bbb", startAt, endAt, "ef0d2079-e9a2-4810-8cae-eb6729c50580", 10000, false, //nolint:lll
	)
	s.Require().NoError(err)

	time.Sleep(10 * time.Second)

	query := `
			SELECT is_notified
			FROM events
			LIMIT 1
		`

	var isNotified bool
	err = s.pgConn.QueryRowContext(s.ctx, query).Scan(&isNotified)
	s.Require().NoError(err)

	s.Require().True(isNotified)
}

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(EventSuite))
}

func (s *EventSuite) setupPG() {
	pgAddr := os.Getenv("POSTGRES_ADDRESS")
	s.Require().NotEmpty(pgAddr)

	pgUser := os.Getenv("POSTGRES_USER")
	s.Require().NotEmpty(pgUser)

	pgPass := os.Getenv("POSTGRES_PASSWORD")
	s.Require().NotEmpty(pgPass)

	pgDB := os.Getenv("POSTGRES_DB")
	s.Require().NotEmpty(pgDB)

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=disable",
		pgUser,
		pgPass,
		pgAddr,
		pgDB,
	)

	pgConn, err := sqlx.Connect("postgres", dsn)

	s.Require().NoError(err)

	s.pgConn = pgConn
}

func (s *EventSuite) setupGRPC() {
	grpcAddr := os.Getenv("GRPC_ADDRESS")
	s.Require().NotEmpty(grpcAddr)

	grpcConn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.grpcConn = grpcConn
	s.eventClient = event.NewEventServiceClient(grpcConn)
}

func (s *EventSuite) cleanTables(tables ...string) {
	_, err := s.pgConn.ExecContext(
		s.ctx,
		fmt.Sprintf("TRUNCATE TABLE %s CASCADE", strings.Join(tables, ", ")),
	)
	s.Require().NoError(err)
}
