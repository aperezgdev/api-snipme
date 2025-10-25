-- name: FindLinkCounterViewCounter :one
SELECT
    id,
    link_id,
    country_code,
    total_views,
    unique_visitors,
    created_on
FROM link_country_view_counter
WHERE link_id = $1 AND country_code = $2;

-- name: SaveLinkCounterViewCounter :exec
INSERT INTO link_country_view_counter (
    id,
    link_id,
    country_code,
    total_views,
    unique_visitors,
    created_on
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: UpdateLinkCounterViewCounter :exec
UPDATE link_country_view_counter
SET
    total_views = $1,
    unique_visitors = $2
WHERE
    link_id = $3 AND country_code = $4;

-- name: RemoveLinkCounterViewCounterByLink :exec
DELETE FROM link_country_view_counter
WHERE link_id = $1;
