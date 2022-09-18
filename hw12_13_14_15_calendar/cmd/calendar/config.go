package main

import (
	"os"

	env "github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

const defaultPath = ".env"

// Config предоставляет настройки приложения.
type Config struct {
	// Log параметры логирования.
	Log LogConfig

	// REST параметры для REST сервера.
	REST RESTConfig

	// PostgreSQL параметры для подключения к PostgreSQL.
	PostgreSQL PostgreSQLConfig

	// DBDriver декларирует драйвер базы данных.
	DBDriver string `env:"DB_DRIVER,required"`
}

// LogConfig предоставляет настройки логирования.
type LogConfig struct {
	// Level уровень логирования.
	Level string `env:"LOG_LEVEL" envDefault:"debug"`
}

// RESTConfig предоставляет настройки REST сервера.
type RESTConfig struct {
	// Address адрес REST сервера.
	Address string `env:"REST_ADDRESS" envDefault:":8080"`
}

// PostgreSQLConfig предоставляет настройки подключения к PostgreSQL.
type PostgreSQLConfig struct {
	// Host адрес БД.
	Host string `env:"POSTGRES_HOST"`

	// Port порт для подключения к БД.
	Port int `env:"POSTGRES_PORT"`

	// User пользователь БД.
	User string `env:"POSTGRES_USER"`

	// Password пароль для подключения к БД.
	Password string `env:"POSTGRES_PASSWORD"`

	// Database название БД.
	Database string `env:"POSTGRES_DB"`
}

// NewConfig создает новый конфиг.
func NewConfig() *Config {
	return &Config{}
}

// LoadConfig создает конфиг на основании переменных окружения.
func loadConfig(path string) (*Config, error) {
	returnErrIfFileNotExists := path != ""

	path = pathOrDefault(path)

	err := godotenv.Load(path)
	if err != nil && (!os.IsNotExist(err) || returnErrIfFileNotExists) {
		return nil, errors.Wrap(err, "cannot load config file")
	}

	cfg := NewConfig()

	if err := env.Parse(cfg); err != nil {
		return nil, errors.Wrap(err, "cannot parse config")
	}

	return cfg, nil
}

// pathOrDefault возвращает путь, который передан, или путь по умолчанию.
func pathOrDefault(path string) string {
	if path == "" {
		return defaultPath
	}
	return path
}
