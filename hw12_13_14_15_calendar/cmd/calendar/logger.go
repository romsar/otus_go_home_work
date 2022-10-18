package main

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// initLogger инициализирует логгер.
func initLogger() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

// configureLogger настраивает логгер.
func configureLogger(config LogConfig) error {
	if config.Level != "" {
		lvl, err := zerolog.ParseLevel(config.Level)
		if err != nil {
			return errors.Wrap(err, "configure log level")
		}

		zerolog.SetGlobalLevel(lvl)
	}

	return nil
}
