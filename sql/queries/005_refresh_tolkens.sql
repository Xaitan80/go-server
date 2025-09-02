-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at, revoked_at, created_at, updated_at)
VALUES ($1, $2, $3, NULL, NOW(), NOW())
RETURNING *;

-- name: GetRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = $2,
    updated_at = $3
WHERE token = $1;

-- name: DeleteAllRefreshTokens :exec
DELETE FROM refresh_tokens;
