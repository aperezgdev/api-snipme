-- name: SaveClient :exec
INSERT INTO client (id, name, email, created_on)
VALUES ($1, $2, $3, $4);

-- name: RemoveClient :exec
DELETE FROM client
WHERE id = $1;
