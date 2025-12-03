-- name: FindClientByID :one
SELECT id, name, email, user_id, created_on
FROM client
WHERE id = $1;

-- name: SaveClient :exec
INSERT INTO client (id, name, email, user_id, created_on)
VALUES ($1, $2, $3, $4, $5);

-- name: RemoveClient :exec
DELETE FROM client
WHERE id = $1;
