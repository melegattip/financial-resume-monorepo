#!/bin/bash
# Financial Resume Monorepo - Build all services

set -e

echo "🔨 Building all services..."

echo ""
echo "📦 Building go-shared..."
cd packages/go-shared && go build ./... && cd ../..

echo ""
echo "📦 Building api-gateway..."
cd apps/api-gateway && go build -o ../../bin/api-gateway ./cmd/api && cd ../..

echo ""
echo "📦 Building users-service..."
cd apps/users-service && go build -o ../../bin/users-service ./cmd/api && cd ../..

echo ""
echo "📦 Building ai-service..."
cd apps/ai-service && go build -o ../../bin/ai-service ./cmd/api && cd ../..

echo ""
echo "📦 Building gamification-service..."
cd apps/gamification-service && go build -o ../../bin/gamification-service ./cmd/api && cd ../..

echo ""
echo "📦 Building frontend..."
cd apps/frontend && npm run build && cd ../..

echo ""
echo "✅ All services built successfully!"
echo "Binaries available in ./bin/"
