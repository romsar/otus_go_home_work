package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"

	"github.com/RomanSarvarov/otus_go_home_work/calendar/config"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/inmem"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/kafka"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/pkg/closer"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/pkg/logging"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/postgres"
	"github.com/RomanSarvarov/otus_go_home_work/calendar/scheduler"
)

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

	var repo scheduler.Repository

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

		repo = r
	default:
		return fmt.Errorf("database driver `%s` not found", cfg.DBDriver)
	}

	w := kafka.NewWriter(&kafka.WriterConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.SenderTopic,
	})
	closer.Add(func() error {
		return w.Close()
	})

	sch := scheduler.New(repo, w, scheduler.Config{
		Interval:        cfg.Scheduler.Interval,
		EventLifeInDays: cfg.Scheduler.EventLifeInDays,
	})

	errgrp.Go(func() error {
		log.
			Debug().
			Msgf("start scheduler")

		return sch.Start(ctx)
	})

	<-ctx.Done()

	log.
		Debug().
		Msg("stopping application")

	closer.CloseAll()

	if err := errgrp.Wait(); err != nil && !errors.Is(err, context.Canceled) {
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
