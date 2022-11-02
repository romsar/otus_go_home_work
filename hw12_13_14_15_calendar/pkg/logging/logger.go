package logging

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Level = zerolog.Level

// InitLogger инициализирует логгер.
func InitLogger() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

// Config предоставляет настройки логирования.
type Config struct {
	Level string
}

// Configure настраивает логгер.
func Configure(config Config) error {
	if config.Level != "" {
		lvl, err := zerolog.ParseLevel(config.Level)
		if err != nil {
			return errors.Wrap(err, "configure log level")
		}

		zerolog.SetGlobalLevel(lvl)
	}

	return nil
}
