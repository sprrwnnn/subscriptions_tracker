#!/usr/bin/env sh
set -eu

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-subscriptions}"
DB_SSLMODE="${DB_SSLMODE:-disable}"
MIGRATION_DIRECTION="${1:-up}"

DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

migrate -path migrations -database "$DATABASE_URL" "$MIGRATION_DIRECTION"
