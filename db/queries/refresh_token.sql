-- name: CreateRefreshToken :one
INSERT INTO refresh_token (id, user_id, token, expires_at, created_on)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_token
WHERE token = $1;

-- name: GetRefreshTokensByUserId :many
SELECT * FROM refresh_token
WHERE user_id = $1
ORDER BY created_on DESC;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_token
WHERE token = $1;

-- name: DeleteRefreshTokensByUserId :exec
DELETE FROM refresh_token
WHERE user_id = $1;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_token
WHERE expires_at < NOW();
