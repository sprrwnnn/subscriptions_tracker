# Subscription Tracker

REST API сервис для учета онлайн-подписок. Приложение позволяет создавать, просматривать, обновлять и удалять подписки, а также считать суммарную стоимость подписок за выбранный период.

Проект оформлен как backend/devops portfolio project: есть Docker, PostgreSQL, миграции, Swagger-документация, структурированные JSON-логи, health checks, Bash-скрипты, GitLab CI и пример деплоя через Ansible.

## Стек

- Go 1.26
- Gin
- GORM
- PostgreSQL
- SQL migrations
- Docker, Docker Compose
- Swagger UI
- Logrus JSON logging
- Bash
- GitLab CI
- Ansible

## Возможности

- CRUD для подписок
- Расчет общей стоимости подписок за период
- Фильтрация списка по `user_id`
- Фильтрация расчета стоимости по `user_id` и `service_name`
- Формат дат подписок `MM-YYYY`, например `07-2025`
- PostgreSQL-хранилище с индексами и проверками данных
- Swagger UI для просмотра API
- `/health` для проверки HTTP-сервера
- `/ready` для проверки готовности приложения и подключения к БД
- JSON-логирование запросов с `request_id`
- Docker image с запуском не от root-пользователя
- GitLab CI pipeline для проверки, сборки и публикации Docker image
- Ansible playbook для деплоя на Linux-сервер

## Архитектура

```text
cmd/main.go
  точка входа: конфиг, БД, router, middleware, health checks

internal/config
  чтение env-переменных и сборка DSN для PostgreSQL

internal/handler
  HTTP handlers, обработка входных запросов и ответов API

internal/service
  бизнес-логика подписок и расчета стоимости

internal/repository
  работа с PostgreSQL через GORM

internal/models
  модели БД, request/response DTO и парсинг дат

internal/middleware
  middleware для структурированного логирования запросов

migrations
  SQL-миграции схемы БД

scripts
  Bash-скрипты для локальной разработки и проверок

ansible
  пример автоматизированного деплоя через Docker Compose
```

## Быстрый запуск через Docker

Скопируйте пример переменных окружения:

```bash
cp .env.example .env
```

Запустите приложение и PostgreSQL:

```bash
make docker-up
```

Или напрямую через Docker Compose:

```bash
docker compose up --build
```

API будет доступно по адресу:

```text
http://localhost:8080/api/v1
```

Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

Если порт `8080` занят:

```bash
HOST_PORT=18080 docker compose up --build
```

## Конфигурация

Переменные окружения из `.env.example`:

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

При запуске в Docker Compose приложение подключается к PostgreSQL по имени сервиса `postgres`.

## Health Checks

Проверка, что HTTP-сервер запущен:

```bash
curl http://localhost:8080/health
```

Ответ:

```json
{
  "status": "ok"
}
```

Проверка готовности приложения и БД:

```bash
curl http://localhost:8080/ready
```

Ответ при доступной БД:

```json
{
  "status": "ready"
}
```

Если PostgreSQL недоступен, `/ready` вернет HTTP `503`.

## Локальная разработка

Установите зависимости:

```bash
go mod tidy
```

Запуск без Docker:

```bash
make run
```

Проверки:

```bash
make fmt
make vet
make test
make check
```

Полезные команды:

```bash
make build
make docker-build
make docker-logs
make docker-ps
make docker-down
make clean
```

Миграции:

```bash
make migrate-up
make migrate-down
```

## Bash-скрипты

В проекте есть вспомогательные скрипты:

```text
scripts/check.sh
  проверяет gofmt, запускает go vet и go test ./...

scripts/dev.sh
  создает .env из .env.example при необходимости и запускает docker compose

scripts/run-migrations.sh
  запускает SQL-миграции через утилиту migrate

scripts/wait-for-db.sh
  ждет доступности PostgreSQL по host/port
```

Примеры:

```bash
./scripts/check.sh
./scripts/dev.sh
./scripts/run-migrations.sh up
DB_HOST=localhost DB_PORT=5432 ./scripts/wait-for-db.sh
```

## API

Базовый путь:

```text
/api/v1
```

| Method | Endpoint | Описание |
| --- | --- | --- |
| `POST` | `/subscriptions` | Создать подписку |
| `GET` | `/subscriptions` | Получить список подписок |
| `GET` | `/subscriptions/:id` | Получить подписку по ID |
| `PUT` | `/subscriptions/:id` | Обновить подписку |
| `DELETE` | `/subscriptions/:id` | Удалить подписку |
| `POST` | `/calculate-cost` | Рассчитать стоимость подписок за период |

### Создать подписку

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

С полем окончания подписки:

```json
{
  "service_name": "Netflix",
  "price": 799,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "01-2025",
  "end_date": "12-2025"
}
```

### Получить список подписок

```bash
curl "http://localhost:8080/api/v1/subscriptions?page=1&page_size=20"
```

С фильтром по пользователю:

```bash
curl "http://localhost:8080/api/v1/subscriptions?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&page=1&page_size=20"
```

### Получить подписку по ID

```bash
curl http://localhost:8080/api/v1/subscriptions/<subscription_id>
```

### Обновить подписку

```bash
curl -X PUT http://localhost:8080/api/v1/subscriptions/<subscription_id> \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 500,
    "start_date": "07-2025"
  }'
```

### Удалить подписку

```bash
curl -X DELETE http://localhost:8080/api/v1/subscriptions/<subscription_id>
```

### Рассчитать стоимость за период

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

Пример ответа:

```json
{
  "total_cost": 2400
}
```

## Docker

В проекте используется multi-stage Dockerfile:

- на первом этапе собирается Go binary;
- на втором этапе запускается минимальный Alpine image;
- приложение работает от отдельного non-root пользователя `app`;
- контейнер содержит `HEALTHCHECK` на `/health`.

Проверить итоговый Docker Compose config:

```bash
docker compose config
```

## GitLab CI

Pipeline описан в `.gitlab-ci.yml`.

Стадии:

- `test`: установка зависимостей, `gofmt`, `go vet`, `go test ./...`
- `build`: сборка бинарного файла `bin/subscription-tracker`
- `package`: сборка и публикация Docker image в GitLab Container Registry

Публикация Docker image выполняется для default branch и Git tags.

## Ansible Deploy

В папке `ansible/` лежит пример деплоя приложения на Linux-сервер с Docker Compose.

Подготовьте inventory:

```bash
cp ansible/inventory.example.ini ansible/inventory.ini
```

Отредактируйте `ansible/inventory.ini` и переменные в `ansible/group_vars/all.example.yml`, затем запустите:

```bash
ansible-playbook -i ansible/inventory.ini ansible/playbook.yml
```

Playbook:

- устанавливает Docker и Docker Compose plugin на Debian-based сервер;
- создает директорию `/opt/subscription-tracker`;
- копирует файлы проекта;
- рендерит `.env`;
- запускает сервис через `docker compose up -d --build`.

## Проверка проекта

Перед коммитом удобно выполнить:

```bash
make check
docker compose config
```

Для проверки Ansible playbook:

```bash
ANSIBLE_LOCAL_TEMP=./tmp/ansible ansible-playbook --syntax-check -i ansible/inventory.example.ini ansible/playbook.yml
```
