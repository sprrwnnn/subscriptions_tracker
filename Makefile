APP_NAME=subscription-tracker
HOST_PORT ?= 8080

.PHONY: run build fmt vet test tidy check docker-build docker-up docker-down docker-logs docker-ps migrate-up migrate-down swagger health ready clean

run:
	go run ./cmd

build:
	go build -o bin/$(APP_NAME) ./cmd

fmt:
	find . -name '*.go' -not -path './.git/*' -not -path './.gocache/*' -not -path './subscriptions_tracker/*' -exec gofmt -w {} +

vet:
	go vet ./...

test:
	go test ./...

tidy:
	go mod tidy

check:
	./scripts/check.sh

docker-build:
	docker compose build

docker-up:
	docker compose up --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f app

docker-ps:
	docker compose ps

migrate-up:
	./scripts/run-migrations.sh up

migrate-down:
	./scripts/run-migrations.sh down

swagger:
	swag init -g cmd/main.go

health:
	curl -fsS http://localhost:$(HOST_PORT)/health

ready:
	curl -fsS http://localhost:$(HOST_PORT)/ready

clean:
	rm -rf bin/
	docker compose down -v
