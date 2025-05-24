#!/usr/bin/env bash
set -eo pipefail

# Check if psql is installed
if ! command -v psql &> /dev/null; then
  echo >&2 "Error: psql is not installed."
  exit 1
fi

# Check if migrate is installed
if ! command -v migrate &> /dev/null; then
  echo >&2 "Error: golang-migrate is not installed."
  echo >&2 "Install it with:"
  echo >&2 "  brew install golang-migrate"
  echo >&2 "  OR follow instructions: https://github.com/golang-migrate/migrate"
  exit 1
fi

# Set your DB config (use your actual local DB settings)
# Adjust all these variables according to your setup 
DB_USER="${POSTGRES_USER:=postgres}"
DB_PASSWORD="${POSTGRES_PASSWORD:=1111}"
DB_NAME="${POSTGRES_DB:=codeLookout}"
DB_PORT="${POSTGRES_PORT:=5432}"
DB_HOST="${POSTGRES_HOST:=localhost}"
MIGRATIONS_DIR="${MIGRATIONS_DIR:-./migrations}" #Adjust this if according to your setup

# Wait for Postgres to be ready on localhost
export PGPASSWORD="${DB_PASSWORD}"
until psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d "${DB_NAME}" -c '\q' &> /dev/null; do
  >&2 echo "Postgres is still unavailable on ${DB_HOST}:${DB_PORT} - sleeping"
  sleep 1
done

echo "Postgres is up and running on ${DB_HOST}:${DB_PORT}!"

# Run migrations
export DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

migrate -path "${MIGRATIONS_DIR}" -database "${DATABASE_URL}" up

echo "Migrations applied successfully!"
