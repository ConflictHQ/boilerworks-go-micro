#!/usr/bin/env bash
set -euo pipefail

# Boilerworks — Go Micro
# Usage: ./run.sh [command]

COMPOSE_FILE=""

if [ -f "docker-compose.yml" ]; then
    COMPOSE_FILE="docker-compose.yml"
elif [ -f "docker-compose.yaml" ]; then
    COMPOSE_FILE="docker-compose.yaml"
fi

compose() {
    if [ -n "$COMPOSE_FILE" ]; then
        docker compose -f "$COMPOSE_FILE" "$@"
    else
        echo "No docker-compose file found"
        exit 1
    fi
}

case "${1:-help}" in
    up|start)
        compose up -d --build
        echo ""
        echo "Services starting. Check status with: ./run.sh status"
        ;;
    down|stop)
        compose down
        ;;
    restart)
        compose down
        compose up -d --build
        ;;
    status|ps)
        compose ps
        ;;
    logs)
        compose logs -f "${2:-}"
        ;;
    seed)
        echo "Set API_KEY_SEED env var and restart the API service"
        echo "docker compose up -d api"
        ;;
    test)
        go test ./...
        ;;
    lint)
        golangci-lint run ./...
        ;;
    shell)
        compose exec api sh
        ;;
    build)
        go build -o api ./cmd/api
        ;;
    help|*)
        echo "Usage: ./run.sh <command>"
        echo ""
        echo "Commands:"
        echo "  up, start     Start all services"
        echo "  down, stop    Stop all services"
        echo "  restart       Restart all services"
        echo "  status, ps    Show service status"
        echo "  logs [svc]    Tail logs (optionally for one service)"
        echo "  seed          Seed the database"
        echo "  test          Run tests"
        echo "  lint          Run linters"
        echo "  shell         Open a shell in the API container"
        echo "  build         Build the Go binary"
        echo "  help          Show this help"
        ;;
esac
