-include .env

MIGRATE=goose -dir=migrations -allow-missing postgres "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

migrate-status:
	$(MIGRATE) status

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

migrate-redo:
	@read -p  "Are you sure to reapply the last migration? [y/n]" -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy] ]]; \
	then \
		$(MIGRATE) redo; \
	fi

migrate-create:
	@read -p  "What is the name of migration?" NAME; \
	${MIGRATE} create $$NAME sql

build:
	go build ./cmd/calendar

run:
	go run ./cmd/calendar --config=.env

test:
	go test -race ./...

lint:
	golangci-lint run ./...