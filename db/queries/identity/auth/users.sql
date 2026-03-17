-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    base_role,
    status,
    email_verified_at,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: MarkUserEmailVerified :exec
UPDATE users
SET status = 'active',
    email_verified_at = $2,
    updated_at = $2
WHERE id = $1;
