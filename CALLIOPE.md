# Calliope — Boilerworks Go Micro
<!-- Agent shim for https://github.com/calliopeai/calliope-cli -->

Primary conventions doc: [`bootstrap.md`](bootstrap.md)
Context seed: [`memory.md`](memory.md)

Read both before writing any code.

---

## Project-specific notes

- Go 1.25+ microservice: Chi router, sqlc (type-safe SQL), goose migrations, pgx/v5 on Postgres 16. API-only — no frontend, no sessions.
- API-key auth (`X-API-Key`, SHA256-hashed) on all endpoints except `/health`; per-key scopes `events.read`, `events.write`, `keys.manage`, `*`.
- UUID primary keys (never expose internal IDs); soft deletes via `deleted_at` (queries filter automatically).
- All responses wrapped in `ApiResponse{Ok, Data, Message, Errors}`; rate limiting is 60 req/min per IP.
- Never edit `internal/database/queries/` by hand — run `sqlc generate`.
- API is on host port 8000 (compose maps `8000:8080`); the seed admin key is printed once on first boot (`docker compose logs api | grep "Plaintext key"`).
