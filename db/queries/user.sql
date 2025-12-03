-- name: CreateUser :one
INSERT INTO "user" (id, email, oauth_provider, oauth_subject, created_on)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM "user"
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM "user"
WHERE email = $1;

-- name: GetUserByOAuthProviderAndSubject :one
SELECT * FROM "user"
WHERE oauth_provider = $1 AND oauth_subject = $2;

-- name: UpdateUser :one
UPDATE "user"
SET email = $2,
    oauth_provider = $3,
    oauth_subject = $4
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE id = $1;
