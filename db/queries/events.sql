-- name: CreateEvent :one
INSERT INTO events (type, payload) VALUES ($1, $2) RETURNING *;

-- name: GetEvent :one
SELECT * FROM events WHERE id = $1 AND deleted_at IS NULL;

-- name: ListEvents :many
SELECT * FROM events WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListEventsByType :many
SELECT * FROM events WHERE type = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: SoftDeleteEvent :exec
UPDATE events SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL;

-- name: CountEvents :one
SELECT COUNT(*) FROM events WHERE deleted_at IS NULL;
