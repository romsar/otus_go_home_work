package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"

	"github.com/RomanSarvarov/otus_go_home_work/calendar/api/rest"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/inmem"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/pkg/closer"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/postgres"
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

	var restModel rest.Model

	switch config.DBDriver {
	case inmem.Key:
		restModel = inmem.New()
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

		restModel = repo
	default:
		return fmt.Errorf("database driver `%s` not found", config.DBDriver)
	}

	srv := newRESTServer(config.REST, restModel)

	closer.Add(func() error {
		log.
			Debug().
			Msgf("terminating REST server")

		if err := srv.Close(ctx); err != nil && !errors.Is(err, context.Canceled) {
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		log.
			Debug().
			Msgf("starting REST server on: `%s`", config.REST.Address)

		err := srv.Listen(rest.LoggingMiddleware)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
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

// newRESTServer создает REST сервер.
func newRESTServer(config RESTConfig, model rest.Model) rest.Server {
	serverCfg := rest.Config{
		Address: config.Address,
	}

	return rest.New(serverCfg, model)
}

// parseFlags возвращает флаги запуска.
func parseFlags() string {
	configPath := flag.StringP("config", "C", "", "Path to configuration file")

	flag.Parse()

	return *configPath
}
