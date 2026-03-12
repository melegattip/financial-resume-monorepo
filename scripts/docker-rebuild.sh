#!/bin/bash
# Financial Resume - Docker Rebuild
# Run this after every new implementation to test changes locally with Docker.

set -e

echo "=== Docker Rebuild: Financial Resume ==="
echo ""

# Parse optional flags
SERVICES=""
NO_CACHE=""

for arg in "$@"; do
  case $arg in
    --no-cache) NO_CACHE="--no-cache" ;;
    backend|frontend) SERVICES="$SERVICES $arg" ;;
    *) echo "Unknown arg: $arg" ;;
  esac
done

# Default: rebuild backend + frontend (not postgres, to preserve data)
if [ -z "$SERVICES" ]; then
  SERVICES="backend frontend"
fi

echo "Services to rebuild: $SERVICES"
[ -n "$NO_CACHE" ] && echo "Mode: --no-cache"
echo ""

# Stop target services (keep postgres running to preserve data)
echo "[1/4] Stopping services..."
docker compose stop $SERVICES

# Remove old containers
echo "[2/4] Removing old containers..."
docker compose rm -f $SERVICES

# Rebuild images
echo "[3/4] Building images..."
docker compose build $NO_CACHE $SERVICES

# Start everything (postgres stays healthy via depends_on)
echo "[4/4] Starting services..."
docker compose up -d $SERVICES

echo ""
echo "=== Done! ==="
echo ""
echo "  Backend:  http://localhost:8080"
echo "  Frontend: http://localhost:3000"
echo ""
echo "Logs:"
echo "  docker compose logs -f backend"
echo "  docker compose logs -f frontend"
