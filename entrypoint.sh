#!/bin/sh
set -e

echo "Running migrations..."
./goose -dir ./migrations postgres "$DATABASE_URL" up

echo "Starting API..."
exec ./jarvis-api
