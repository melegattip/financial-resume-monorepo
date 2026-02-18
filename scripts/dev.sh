#!/bin/bash
# Financial Resume - Local Development

set -e

echo "Starting local dev environment..."

# Start postgres
docker compose up -d postgres

echo "Waiting for postgres..."
sleep 3

echo "Postgres ready!"
echo ""
echo "Start the backend:"
echo "  cd apps/monolith && go run ./cmd/server/"
echo ""
echo "Start the frontend:"
echo "  cd apps/frontend && npm start"
