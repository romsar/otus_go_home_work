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
	go build ./cmd/calendar_scheduler
	go build ./cmd/calendar_sender

run:
	go run ./cmd/calendar --config=.env

run-scheduler:
	go run ./cmd/calendar_scheduler --config=.env

run-sender:
	go run ./cmd/calendar_sender --config=.env

up:
	cp -n .env.example .env || true
	docker-compose -f ./deployments/docker-compose.yml --env-file .env -p calendar up -d --build

down:
	docker-compose -f ./deployments/docker-compose.yml --env-file .env -p calendar stop

test:
	go test -race -count 10 ./...

integration-tests:
	sh ./scripts/run-integration-test.sh

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

generate:
	go generate ./...

.PHONY: build run build-img run-img version test lint proto generate