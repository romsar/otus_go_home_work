package main

import (
	"context"
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
	"github.com/RomanSarvarov/otus_go_home_work/calendar/sender"
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

	var model sender.Model

	switch cfg.DBDriver {
	case inmem.Key:
		model = inmem.New()
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

		model = repo
	default:
		return fmt.Errorf("database driver `%s` not found", cfg.DBDriver)
	}
	fmt.Println(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.GroupID,
		Topic:   cfg.Kafka.SenderTopic,
	})

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.GroupID,
		Topic:   cfg.Kafka.SenderTopic,
	})
	closer.Add(func() error {
		return r.Close()
	})

	s := sender.New(model, r, sender.Config{
		Threads: cfg.Sender.Threads,
	})

	errgrp.Go(func() error {
		log.
			Debug().
			Msgf("start sender")

		return s.Start(ctx)
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
