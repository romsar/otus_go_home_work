package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres support
	"github.com/pkg/errors"
	goose "github.com/pressly/goose/v3"
)

// Key обозначает ключ postgres БД драйвера.
const Key = "postgres"

// Config содрежит настройки подключения к БД.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// Repository является абстракцией к БД PostgreSQL.
type Repository struct {
	db *sqlx.DB
}

// Open открывает соединение к БД.
func Open(cfg Config) (*Repository, error) {
	db, err := sqlx.Connect("postgres", dsn(cfg))
	if err != nil {
		return nil, errors.Wrap(err, "open postgreSQL connection")
	}

	return &Repository{db: db}, nil
}

// Close закрывает подключение к БД.
func (repo *Repository) Close() error {
	return repo.db.Close()
}

// Up выполняет миграции.
func (repo *Repository) Up(dir string) error {
	return goose.Up(repo.db.DB, dir)
}

// dsn формирует DSN строку из конфига.
func dsn(cfg Config) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
}
