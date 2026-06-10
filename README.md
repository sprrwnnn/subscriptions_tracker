# Subscription Tracker

REST-сервис для агрегации пользовательских онлайн-подписок

## Возможности

- CRUD-операции с подписками
- Расчёт общей стоимости подписок за выбранный период
- Опциональная фильтрация по ID пользователя и названию сервиса
- Хранение в PostgreSQL с SQL-миграциями
- Структурированное JSON-логирование
- Swagger UI.
- Запуск через Docker Compose

## Конфигурация

Скопируйте .env.example в .env для локального запуска:

```bash
cp .env.example .env
```

Доступные переменные:

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

## Запуск с Docker Compose

```bash
docker compose up --build
```

Если порт 8080 уже занят, выберите другой порт на хосте:

```bash
HOST_PORT=18080 docker compose up --build
```

API будет доступно по адресу:

```text
http://localhost:8080/api/v1
```

Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

## Примеры API

Создание подписки:

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

Список подписок:

```bash
curl "http://localhost:8080/api/v1/subscriptions?page=1&page_size=20"
```

Расчет общей стоимости:

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

## Локальная разработка

```bash
go mod tidy
go test ./...
go run ./cmd
```
