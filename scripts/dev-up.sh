#!/bin/bash
# Financial Resume Monorepo - Development Environment

set -e

echo "🚀 Starting Financial Resume development environment..."

# Start databases
echo "📦 Starting databases..."
docker-compose -f infrastructure/docker/docker-compose.yml up -d postgres-main postgres-users redis

# Wait for databases to be ready
echo "⏳ Waiting for databases..."
sleep 5

# Run migrations (if applicable)
# echo "🔄 Running migrations..."
# ./scripts/migrate.sh

echo "✅ Development environment ready!"
echo ""
echo "Start services individually:"
echo "  cd apps/api-gateway && go run cmd/api/main.go"
echo "  cd apps/users-service && go run cmd/api/main.go"
echo "  cd apps/ai-service && go run cmd/api/main.go"
echo "  cd apps/gamification-service && go run cmd/api/main.go"
echo "  cd apps/frontend && npm start"
