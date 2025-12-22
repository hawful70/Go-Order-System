COMPOSE := docker compose -f infrastructure/docker-compose.dev.yml
GOCACHE_DIR := $(PWD)/.cache/go

.PHONY: help run test docker-up docker-down docker-logs

help:
	@echo "Available targets:"
	@echo "  run          - Run the email service locally"
	@echo "  test         - Execute Go tests with a local cache dir"
	@echo "  docker-up    - Start the email service container"
	@echo "  docker-down  - Stop the email service container"
	@echo "  docker-logs  - Tail logs from the email service container"

run:
	go run ./cmd/email-service

test:
	mkdir -p $(GOCACHE_DIR)
	GOCACHE=$(GOCACHE_DIR) go test ./...

docker-up:
	$(COMPOSE) up -d --build

docker-down:
	$(COMPOSE) down

docker-logs:
	$(COMPOSE) logs -f email-service
