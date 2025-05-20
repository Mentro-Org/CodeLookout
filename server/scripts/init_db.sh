#!/usr/bin/env bash
set -x
set -eo pipefail

if ! [ -x "$(command -v psql)" ]; then
  echo >&2 "Error: psql is not installed."
  exit 1
fi

if ! [ -x "$(command -v migrate)" ]; then
  echo >&2 "Error: golang-migrate is not installed."
  echo >&2 "Install it with:"
  echo >&2 "  brew install golang-migrate"  # macOS
  echo >&2 "  OR follow instructions: https://github.com/golang-migrate/migrate"
  exit 1
fi

DB_USER="${POSTGRES_USER:=postgres}"
DB_PASSWORD="${POSTGRES_PASSWORD:=1234}"
DB_NAME="${POSTGRES_DB:=newsletter}"
DB_PORT="${POSTGRES_PORT:=5432}"

# Allow skipping Docker if Postgres is already running
if [[ -z "${SKIP_DOCKER}" ]]; then
  docker run \
    -e POSTGRES_USER="${DB_USER}" \
    -e POSTGRES_PASSWORD="${DB_PASSWORD}" \
    -e POSTGRES_DB="${DB_NAME}" \
    -p "${DB_PORT}":5432 \
    -d postgres \
    postgres -N 1000
fi

# Wait for Postgres to be ready
export PGPASSWORD="${DB_PASSWORD}"
until psql -h "localhost" -U "${DB_USER}" -p "${DB_PORT}" -d "postgres" -c '\q'; do
  >&2 echo "Postgres is still unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up and running on port ${DB_PORT}!"

# Run migrations
MIGRATIONS_DIR=/home/ananya/CodeLookout/server/migrations  # Adjust this if needed
export DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable"

migrate -path "${MIGRATIONS_DIR}" -database "${DATABASE_URL}" up

>&2 echo "Postgres has been migrated, ready to go!"

exec go run cmd/main.go