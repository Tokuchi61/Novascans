SQLC=sqlc
COMPOSE=docker compose
COMPOSE_PROJECT_NAME=novascans
GO_TEST_IMAGE=golang:1.26.1-alpine
TEST_DB_NAME=novascans_test
TEST_DB_DSN=postgres://postgres:postgres@postgres:5432/$(TEST_DB_NAME)?sslmode=disable

.PHONY: fmt test test-integration test-db-ensure dev-up dev-down sqlc migrate-up migrate-down migrate-status

fmt:
	go fmt ./...

test:
	go test ./...

test-integration: test-db-ensure
	docker run --rm --network $(COMPOSE_PROJECT_NAME)_default -v "$(CURDIR):/workspace" -w /workspace -e NOVASCANS_TEST_DATABASE_URL=$(TEST_DB_DSN) -e GOCACHE=/workspace/.cache/go-build $(GO_TEST_IMAGE) go test -tags=integration ./...

test-db-ensure:
	$(COMPOSE) up -d postgres
	$(COMPOSE) exec -T postgres sh /docker-entrypoint-initdb.d/01-create-test-db.sh

dev-up:
	$(COMPOSE) up --build -d

dev-down:
	$(COMPOSE) down

sqlc:
	$(SQLC) generate

migrate-up:
	$(COMPOSE) exec -T api /app/migrate up

migrate-down:
	$(COMPOSE) exec -T api /app/migrate down

migrate-status:
	$(COMPOSE) exec -T api /app/migrate status
