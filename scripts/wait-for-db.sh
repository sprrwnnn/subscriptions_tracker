#!/usr/bin/env sh
set -eu

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-30}"

start_time="$(date +%s)"

while ! nc -z "$DB_HOST" "$DB_PORT"; do
  now="$(date +%s)"
  elapsed="$((now - start_time))"

  if [ "$elapsed" -ge "$TIMEOUT_SECONDS" ]; then
    echo "Database ${DB_HOST}:${DB_PORT} is not available after ${TIMEOUT_SECONDS}s"
    exit 1
  fi

  echo "Waiting for database ${DB_HOST}:${DB_PORT}..."
  sleep 2
done

echo "Database ${DB_HOST}:${DB_PORT} is available"
