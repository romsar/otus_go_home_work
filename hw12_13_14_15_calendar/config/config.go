package config

import (
	"os"
	"time"

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

	// GRPC параметры для GRPC сервера.
	GRPC GRPCConfig

	// PostgreSQL параметры для подключения к PostgreSQL.
	PostgreSQL PostgreSQLConfig

	// Kafka настройки работы с Kafka.
	Kafka KafkaConfig

	// Scheduler настройки планировщика.
	Scheduler SchedulerConfig

	// Sender настройки отправителя.
	Sender SenderConfig

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

// GRPCConfig предоставляет настройки GRPC сервера.
type GRPCConfig struct {
	// Address адрес GRPC сервера.
	Address string `env:"GRPC_ADDRESS" envDefault:":8081"`
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

// KafkaConfig предоставляет настройки работы с Kafka.
type KafkaConfig struct {
	// GroupID адреса брокеров.
	Brokers []string `env:"KAFKA_BROKERS" envDefault:"kafka:9092" envSeparator:","`

	// GroupID идентификатор группы.
	GroupID string `env:"KAFKA_GROUP_ID" envDefault:"calendar"`

	// SenderTopic название топика для планировщика.
	SenderTopic string `env:"KAFKA_SENDER_TOPIC" envDefault:"calendar-sender-topic"`
}

// SchedulerConfig предоставляет настройки планировщика.
type SchedulerConfig struct {
	// Interval интервал работы планировщика.
	Interval time.Duration `env:"SCHEDULER_INTERVAL" envDefault:"1m"`

	// EventLifeInDays количество дней, после истечения которых удалять событие.
	EventLifeInDays uint `env:"SCHEDULER_EVENT_LIFE_IN_DAYS" envDefault:"365"`
}

// SenderConfig предоставляет настройки отправителя.
type SenderConfig struct {
	// Threads количество потоков (консьюмеров).
	Threads int `env:"SENDER_THREADS" envDefault:"3"`
}

// NewConfig создает новый конфиг.
func NewConfig() *Config {
	return &Config{}
}

// Load создает конфиг на основании переменных окружения.
func Load(path string) (*Config, error) {
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
