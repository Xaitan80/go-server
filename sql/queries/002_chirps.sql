-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id, author_id)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2, $2)
RETURNING *;

-- name: GetChirpByID :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: GetAllChirps :many
SELECT *
FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpsByAuthorID :many
SELECT *
FROM chirps
WHERE author_id = $1
ORDER BY created_at ASC;

-- name: DeleteAllChirps :exec
DELETE FROM chirps;
