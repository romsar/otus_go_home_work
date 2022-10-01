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
	"github.com/philip-bui/grpc-zerolog"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"

	"google.golang.org/grpc"

	"github.com/RomanSarvarov/otus_go_home_work/calendar"
	grpcapi "github.com/RomanSarvarov/otus_go_home_work/calendar/api/grpc"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/inmem"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/pkg/closer"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/postgres"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/proto/event"
)

// migrationsDir определяет местонахождение миграций.
const migrationsDir = "migrations"

func main() {
	initLogger()

	log.Info().Msg("start")

	cfgPath := parseFlags()

	log.
		Debug().
		Str("cfg path", cfgPath).
		Msg("flags parsed")

	config, err := loadConfig(cfgPath)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.
		Debug().
		Interface("config", config).
		Msg("config loaded")

	if err := configureLogger(config.Log); err != nil {
		log.Fatal().Err(err).Send()
	}

	if err := run(config); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// run запускате приложение.
func run(config *Config) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	defer closer.CloseAll()

	errgrp, ctx := errgroup.WithContext(ctx)

	log.
		Debug().
		Msg("start application")

	var model calendar.Model

	switch config.DBDriver {
	case inmem.Key:
		model = inmem.New()
	case postgres.Key:
		log.
			Debug().
			Msg("connecting to postgres")

		dbCfg := postgres.Config{
			Host:     config.PostgreSQL.Host,
			Port:     config.PostgreSQL.Port,
			User:     config.PostgreSQL.User,
			Password: config.PostgreSQL.Password,
			Database: config.PostgreSQL.Database,
		}

		repo, err := postgres.Open(dbCfg)
		if err != nil {
			return err
		}

		closer.Add(func() error {
			log.
				Debug().
				Msgf("terminating postgres connection")

			if err := repo.Close(); err != nil {
				return err
			}

			return nil
		})

		log.
			Debug().
			Msgf("run postgres migrations")

		if err := repo.Up(migrationsDir); err != nil {
			return err
		}

		model = repo
	default:
		return fmt.Errorf("database driver `%s` not found", config.DBDriver)
	}

	// Start REST.
	mux := runtime.NewServeMux()
	restSrv := &http.Server{
		Addr:    config.REST.Address,
		Handler: mux,
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
			Msgf("starting REST server on: `%s`", config.REST.Address)

		err := event.RegisterEventServiceHandlerServer(context.Background(), mux, grpcapi.New(model))
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
	lis, err := net.Listen("tcp", config.GRPC.Address)
	if err != nil {
		return errors.Wrap(err, "listen tcp for grpc")
	}

	grpcSrv := grpc.NewServer(
		zerolog.UnaryInterceptor(),
	)

	event.RegisterEventServiceServer(grpcSrv, grpcapi.New(model))

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
			Msgf("starting GRPC server on: `%s`", config.GRPC.Address)

		err := grpcSrv.Serve(lis)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return err
		}
		return nil
	})

	<-ctx.Done()

	log.
		Debug().
		Msg("stoping application")

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
