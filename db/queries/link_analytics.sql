-- name: SaveLinkAnalytics :exec
INSERT INTO link_analytics (
    id,
    link_id,
    total_views,
    unique_visitors,
    created_on
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: UpdateLinkAnalytics :exec
UPDATE link_analytics
SET
    total_views = $2,
    unique_visitors = $3,
    created_on = COALESCE($4, created_on)
WHERE
    id = $1;

-- name: RemoveLinkAnalyticsByLink :exec
DELETE FROM link_analytics
WHERE link_id = $1;
