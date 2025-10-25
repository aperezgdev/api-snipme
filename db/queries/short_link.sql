-- name: SaveShortLink :exec
INSERT INTO short_link (id, original_route, client_id, code, created_on)
VALUES ($1, $2, $3, $4, $5);

-- name: FindShortLinkByCode :one
SELECT
    id,
    original_route,
    client_id,
    code,
    created_on
FROM short_link
WHERE code = $1;

-- name: RemoveShortLink :exec
DELETE FROM short_link
WHERE id = $1;
