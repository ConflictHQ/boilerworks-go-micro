# Claude -- Boilerworks Go Micro

Primary conventions doc: [`bootstrap.md`](bootstrap.md)

Read it before writing any code.

## Stack

- **Backend**: Go 1.22+ (Chi router)
- **Frontend**: None (API-only microservice)
- **API**: REST with JSON responses
- **Queries**: sqlc (type-safe SQL)
- **Migrations**: goose
- **Auth**: API-key (SHA256 hashed, per-key scopes)

## Quick Reference

| Endpoint | URL |
|----------|-----|
| Health | http://localhost:8080/health |
| Events | http://localhost:8080/events |
| API Keys | http://localhost:8080/api-keys |

## Commands

```bash
make up        # Start Docker services
make down      # Stop services
make build     # Build binary
make test      # Run tests
make lint      # golangci-lint
make logs      # Tail container logs
```

## Structure

```
cmd/api/main.go          — entry point, seed logic
internal/
  config/                — env config loader
  server/                — Chi router setup + route registration
  middleware/             — API-key auth, scope checking, rate limiter
  handler/               — health, event CRUD, API key management
  database/              — pgx pool connection
  database/queries/      — sqlc generated (DO NOT edit)
db/
  migrations/            — goose SQL migrations
  queries/               — sqlc query definitions
  sqlc.yaml              — sqlc config
```

## Rules

- API-key auth on all endpoints except /health
- UUID primary keys, never expose internal IDs
- Soft deletes on business models (deletedAt field)
- Scopes: `events.read`, `events.write`, `keys.manage`, `*` (wildcard)
- All responses wrapped in `ApiResponse{Ok, Data, Message, Errors}`
- Rate limiting: 60 requests per minute per IP
