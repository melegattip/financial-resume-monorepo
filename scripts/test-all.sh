#!/bin/bash
# Financial Resume Monorepo - Run all tests

set -e

echo "🧪 Running all tests..."

echo ""
echo "📦 Testing go-shared..."
cd packages/go-shared && go test ./... && cd ../..

echo ""
echo "📦 Testing api-gateway..."
cd apps/api-gateway && go test ./... && cd ../..

echo ""
echo "📦 Testing users-service..."
cd apps/users-service && go test ./... && cd ../..

echo ""
echo "📦 Testing ai-service..."
cd apps/ai-service && go test ./... && cd ../..

echo ""
echo "📦 Testing gamification-service..."
cd apps/gamification-service && go test ./... && cd ../..

echo ""
echo "📦 Testing frontend..."
cd apps/frontend && npm test -- --watchAll=false && cd ../..

echo ""
echo "✅ All tests passed!"
