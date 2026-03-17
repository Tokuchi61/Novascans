-- name: CreatePasswordCredential :one
INSERT INTO auth_password_credentials (
    user_id,
    password_hash,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetPasswordCredentialByUserID :one
SELECT * FROM auth_password_credentials
WHERE user_id = $1
LIMIT 1;

-- name: GetPasswordCredentialByEmail :one
SELECT
    apc.user_id,
    apc.password_hash,
    apc.created_at,
    apc.updated_at
FROM auth_password_credentials apc
INNER JOIN users u ON u.id = apc.user_id
WHERE u.email = $1
LIMIT 1;

-- name: UpdatePasswordHash :exec
UPDATE auth_password_credentials
SET password_hash = $2,
    updated_at = $3
WHERE user_id = $1;
