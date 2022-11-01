package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpczerolog "github.com/philip-bui/grpc-zerolog"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"

	"google.golang.org/grpc"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
	grpcapi "github.com/RomanSarvarov/otus_go_home_work/calendar/api/grpc"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/api/rest"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/config"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/inmem"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/pkg/closer"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/pkg/logging"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/postgres"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/proto/event"
)

// migrationsDir определяет местонахождение миграций.
const migrationsDir = "migrations"

func main() {
	logging.InitLogger()

	log.Info().Msg("start")

	cfgPath := parseFlags()

	log.
		Debug().
		Str("cfg path", cfgPath).
		Msg("flags parsed")

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.
		Debug().
		Interface("config", cfg).
		Msg("config loaded")

	logConfig := logging.Config{Level: cfg.Log.Level}
	if err := logging.Configure(logConfig); err != nil {
		log.Fatal().Err(err).Send()
	}

	if err := run(cfg); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// run запускает приложение.
func run(cfg *config.Config) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	defer closer.CloseAll()

	errgrp, ctx := errgroup.WithContext(ctx)

	log.
		Debug().
		Msg("start application")

	var repo calendar.Repository

	switch cfg.DBDriver {
	case inmem.Key:
		repo = inmem.New()
	case postgres.Key:
		log.
			Debug().
			Msg("connecting to postgres")

		dbCfg := postgres.Config{
			Host:     cfg.PostgreSQL.Host,
			Port:     cfg.PostgreSQL.Port,
			User:     cfg.PostgreSQL.User,
			Password: cfg.PostgreSQL.Password,
			Database: cfg.PostgreSQL.Database,
		}

		r, err := postgres.Open(dbCfg)
		if err != nil {
			return err
		}

		closer.Add(func() error {
			log.
				Debug().
				Msgf("terminating postgres connection")

			return r.Close()
		})

		log.
			Debug().
			Msgf("run postgres migrations")

		if err := r.Up(migrationsDir); err != nil {
			return err
		}

		repo = r
	default:
		return fmt.Errorf("database driver `%s` not found", cfg.DBDriver)
	}

	// Start REST.
	mux := runtime.NewServeMux()
	restSrv := &http.Server{
		Addr:    cfg.REST.Address,
		Handler: rest.LoggingMiddleware(mux),
	}

	closer.Add(func() error {
		log.
			Debug().
			Msgf("terminating REST server")

		if err := restSrv.Close(); err != nil && !errors.Is(err, context.Canceled) {
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		log.
			Debug().
			Msgf("starting REST server on: `%s`", cfg.REST.Address)

		err := event.RegisterEventServiceHandlerServer(context.Background(), mux, grpcapi.New(repo))
		if err != nil {
			return errors.Wrap(err, "register event service handler server")
		}

		err = restSrv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	// Start GRPC.
	lis, err := net.Listen("tcp", cfg.GRPC.Address)
	if err != nil {
		return errors.Wrap(err, "listen tcp for grpc")
	}

	grpcSrv := grpc.NewServer(
		grpczerolog.UnaryInterceptor(),
	)

	event.RegisterEventServiceServer(grpcSrv, grpcapi.New(repo))

	closer.Add(func() error {
		log.
			Debug().
			Msgf("terminating GRPC server")

		grpcSrv.GracefulStop()

		return nil
	})

	errgrp.Go(func() error {
		log.
			Debug().
			Msgf("starting GRPC server on: `%s`", cfg.GRPC.Address)

		err := grpcSrv.Serve(lis)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return err
		}
		return nil
	})

	<-ctx.Done()

	log.
		Debug().
		Msg("stopping application")

	closer.CloseAll()

	if err := errgrp.Wait(); err != nil {
		return err
	}

	log.
		Debug().
		Msg("application was stopped gracefully")

	return nil
}

// parseFlags возвращает флаги запуска.
func parseFlags() string {
	configPath := flag.StringP("config", "C", "", "Path to configuration file")

	flag.Parse()

	return *configPath
}
