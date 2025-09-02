-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (id, token, user_id, expires_at, revoked_at)
VALUES (gen_random_uuid(), $1, $2, $3, $4)
RETURNING *;

-- name: GetRefreshTokenByToken :one
SELECT *
FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE token = $1;

-- name: DeleteAllRefreshTokens :exec
DELETE FROM refresh_tokens;
