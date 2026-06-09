# Subscription Tracker

REST-service for aggregating user online subscriptions.

## Features

- CRUD operations for subscriptions.
- Total subscription cost calculation for a selected period.
- Optional filtering by user ID and service name.
- PostgreSQL storage with SQL migrations.
- Structured JSON logging.
- Swagger UI.
- Docker Compose startup.

## Configuration

Copy `.env.example` to `.env` for local runs:

```bash
cp .env.example .env
```

Available variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions
DB_SSLMODE=disable
SERVER_PORT=8080
HOST_PORT=8080
LOG_LEVEL=info
```

## Run With Docker Compose

```bash
docker compose up --build
```

If port `8080` is already busy, choose another host port:

```bash
HOST_PORT=18080 docker compose up --build
```

The API will be available at:

```text
http://localhost:8080/api/v1
```

Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

## API Examples

Create a subscription:

```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

List subscriptions:

```bash
curl "http://localhost:8080/api/v1/subscriptions?page=1&page_size=20"
```

Calculate total cost:

```bash
curl -X POST http://localhost:8080/api/v1/calculate-cost \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "service_name": "Yandex Plus",
    "start_date": "07-2025",
    "end_date": "12-2025"
  }'
```

## Local Development

```bash
go mod tidy
go test ./...
go run ./cmd
```
