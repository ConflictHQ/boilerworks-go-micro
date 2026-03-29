# Boilerworks Go Micro

> Lightweight Go microservice with Chi router, sqlc, goose, and API-key auth.
> No frontend, no sessions -- pure API service.

## Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.22+ (Chi router) |
| Queries | sqlc (type-safe SQL) |
| Migrations | goose |
| Database | PostgreSQL 16 |
| Auth | API-key (SHA256, per-key scopes) |
| Linting | golangci-lint |

## Getting Started

```bash
# Start services
docker compose up -d --build

# Get your seed API key (shown once on first boot)
docker compose logs api | grep "Plaintext key"

# Test it
curl http://localhost:8080/health
curl -H "X-API-Key: bw_seed_key_change_me_in_production" http://localhost:8080/events
```

## Endpoints

| Method | Path | Auth | Scope | Description |
|--------|------|------|-------|-------------|
| GET | /health | None | - | Health check |
| POST | /events | API Key | events.write | Create event |
| GET | /events | API Key | events.read | List events |
| GET | /events/{id} | API Key | events.read | Event detail |
| DELETE | /events/{id} | API Key | events.write | Soft delete |
| POST | /api-keys | API Key | keys.manage | Create key |
| GET | /api-keys | API Key | keys.manage | List keys |
| DELETE | /api-keys/{id} | API Key | keys.manage | Revoke key |

## Commands

```bash
make up             # Start Docker services
make down           # Stop services
make build          # Build binary
make test           # Run tests
make lint           # golangci-lint
make migrate-up     # Run migrations
make migrate-down   # Rollback migration
make sqlc-generate  # Regenerate sqlc queries
make logs           # Tail container logs
```

## Documentation

- [bootstrap.md](bootstrap.md) -- Conventions and patterns
- [CLAUDE.md](CLAUDE.md) -- Agent shim

---

Boilerworks is a [Conflict](https://weareconflict.com) brand. CONFLICT is a registered trademark of Conflict LLC.
