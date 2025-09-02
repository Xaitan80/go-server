-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1
RETURNING id, user_id;
