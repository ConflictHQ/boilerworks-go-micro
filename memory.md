# Boilerworks Memory

This file is the **AI context seed** for the Boilerworks Go Micro template. It captures decisions, constraints, and non-obvious facts that are not derivable from reading the code.

For conventions and patterns, see [`bootstrap.md`](bootstrap.md).

---

## Template purpose

Lightweight Go microservice: Chi router, sqlc queries, goose migrations, API-key auth. No frontend, no sessions -- pure JSON API.

---

## Current state

| Component | Version / fact |
|---|---|
| Go | 1.25 (go.mod requires 1.25.5) |
| Router | github.com/go-chi/chi/v5 v5.2.5 |
| Postgres driver | github.com/jackc/pgx/v5 v5.9.1 (pool in `internal/database`) |
| Database | PostgreSQL 16 (compose publishes host port 5432) |
| API | container listens on 8080; compose maps host **8000** -> 8080 |
| Migrations | goose (SQL files in `db/migrations`, run by the `migrate` compose service) |
| Queries | sqlc -- generated code in `internal/database/queries` (never edit by hand) |

---

## Things that bite newcomers

- **The API is on host port 8000, not 8080** -- compose maps `8000:8080`; only in-container URLs use 8080.
- **The seed API key is printed once** on first boot (`docker compose logs api | grep "Plaintext key"`); keys are SHA256-hashed, plaintext is never stored.
- **All endpoints except `/health` require `X-API-Key`** with a matching scope (`events.read`, `events.write`, `keys.manage`, or `*`).
- **Soft deletes only** -- business rows get `deleted_at`; queries filter automatically.
- **UUID primary keys** -- never expose internal IDs.
- **Rate limiting** is 60 requests/minute/IP by default.
