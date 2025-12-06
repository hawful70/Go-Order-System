COMPOSE := docker compose -f infrastructure/docker-compose.yml
GOCACHE_DIR := $(PWD)/.cache/go
PSQL := $(COMPOSE) exec -T postgres env PGPASSWORD=postgres psql -U postgres -d identity

.PHONY: help run test docker-up docker-down docker-logs db-seed

help:
	@echo "Available targets:"
	@echo "  run          - Run the identity service locally"
	@echo "  test         - Execute Go tests with a local cache dir"
	@echo "  docker-up    - Start Postgres, pgAdmin, and the identity service containers"
	@echo "  docker-down  - Stop all containers and remove them"
	@echo "  docker-logs  - Tail logs from the identity service container"
	@echo "  db-seed      - Apply local seed data via psql"

run:
	go run ./cmd/identity-service

test:
	mkdir -p $(GOCACHE_DIR)
	GOCACHE=$(GOCACHE_DIR) go test ./...

docker-up:
	$(COMPOSE) up -d --build

docker-down:
	$(COMPOSE) down

docker-logs:
	$(COMPOSE) logs -f identity-service

db-seed:
	cat migrations/001_seed_users.sql | $(PSQL)
