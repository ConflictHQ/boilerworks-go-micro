# Boilerworks Go Micro -- Bootstrap

> Go 1.22+ microservice with Chi router, sqlc queries, goose migrations,
> and API-key authentication. No frontend, no sessions -- pure API service.

## Architecture

```
Caller (service, cron, webhook sender)
  |
  v (HTTP + X-API-Key header)
  |
Go (Chi router)
  |-- sqlc (Postgres 16)
  +-- JSON API responses
```

## Conventions

### Auth
- All endpoints require `X-API-Key` header except `/health`
- Keys are SHA256-hashed before storage -- plaintext never stored
- Per-key scopes: `events.read`, `events.write`, `keys.manage`, `*`
- `ApiKeyAuth` middleware validates key and stores in context
- `RequireScope` middleware checks scope on individual routes

### Models
- UUID primary keys (`gen_random_uuid()`)
- Snake_case table and column names
- Audit fields: `created_at`, `updated_at`
- Soft deletes: `deleted_at` field, queries filter automatically

### API
- All responses wrapped in `ApiResponse`: `{ ok, data, message, errors }`
- JSON request bodies decoded with `encoding/json`
- Rate limiting: 60 requests per minute per IP (configurable)

### Database
- sqlc generates type-safe Go from SQL queries
- goose manages migrations (SQL format)
- pgx/v5 as the Postgres driver with connection pooling
- Never edit files in `internal/database/queries/` -- run `sqlc generate`

### Docker
- `docker compose up -d --build` starts API + Postgres + migrations
- Goose runs automatically before the API starts
- Seed creates admin key with `['*']` scopes (logged to stdout once)

### Seed API Key
On first boot, check container logs for the plaintext key:
```bash
docker compose logs api | grep "Plaintext key"
```
