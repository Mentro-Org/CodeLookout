#!/usr/bin/env bash
set -eo pipefail

if ! [ -x "$(command -v migrate)" ]; then
  echo >&2 "Error: golang-migrate is not installed."
  echo >&2 "Install it with:"
  echo >&2 "  brew install golang-migrate"  # macOS
  echo >&2 "  OR follow instructions: https://github.com/golang-migrate/migrate"
  exit 1
fi

# Load environment variables from .env file
if [ -f .env ]; then
  set -a
  source .env >/dev/null 2>&1
  set +a
else
  echo ".env file not found."
  exit 1
fi

# Validate DATABASE_URL
if [ -z "${DATABASE_URL}" ]; then
  echo >&2 "Error: DATABASE_URL is not set."
  echo >&2 "Please export DATABASE_URL or define it in a .env file."
  exit 1
fi

# Run migrations
migrate -path "./migrations" -database "${DATABASE_URL}" up

echo >&2 "DB migration successfully!"