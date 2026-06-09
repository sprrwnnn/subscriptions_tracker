APP_NAME=subscription-tracker

.PHONY: run build test tidy docker-build docker-up docker-down migrate-up migrate-down swagger clean

run:
	go run ./cmd

build:
	go build -o bin/$(APP_NAME) ./cmd

test:
	go test ./...

tidy:
	go mod tidy

docker-build:
	docker compose build

docker-up:
	docker compose up --build

docker-down:
	docker compose down

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/subscriptions?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/subscriptions?sslmode=disable" down

swagger:
	swag init -g cmd/main.go

clean:
	rm -rf bin/
	docker compose down -v
