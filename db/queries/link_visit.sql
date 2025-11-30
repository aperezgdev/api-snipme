-- name: SaveLinkVisit :exec
INSERT INTO link_visit (id, link_id, ip, user_agent, created_on)
VALUES ($1, $2, $3, $4, $5);

-- name: FindOldLinkVisits :many
SELECT id, link_id, ip, user_agent, created_on
FROM link_visit
WHERE created_on < NOW() - INTERVAL '15 minutes';

-- name: RemoveLinkVisits :exec
DELETE FROM link_visit
WHERE id = ANY($1::uuid[]);