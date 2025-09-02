-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserFromRefreshToken :one
SELECT 
    u.id AS user_id,
    u.email,
    u.hashed_password,
    u.created_at AS user_created_at,
    u.updated_at AS user_updated_at,
    r.token,
    r.expires_at,
    r.revoked_at
FROM users u
JOIN refresh_tokens r ON u.id = r.user_id
WHERE r.token = $1;

