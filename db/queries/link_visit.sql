-- name: SaveLinkVisit :exec
INSERT INTO link_visit (id, link_id, ip, user_agent, created_on)
VALUES ($1, $2, $3, $4, $5);
