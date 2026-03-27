-- name: GetApiKeyByHash :one
SELECT * FROM api_keys WHERE key_hash = $1 AND is_active = TRUE;

-- name: ListApiKeys :many
SELECT id, name, scopes, is_active, last_used_at, created_at FROM api_keys ORDER BY created_at DESC;

-- name: CreateApiKey :one
INSERT INTO api_keys (name, key_hash, scopes) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateLastUsed :exec
UPDATE api_keys SET last_used_at = NOW() WHERE id = $1;

-- name: RevokeApiKey :exec
UPDATE api_keys SET is_active = FALSE WHERE id = $1;
