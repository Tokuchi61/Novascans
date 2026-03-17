-- name: CreatePasswordResetToken :one
INSERT INTO auth_password_reset_tokens (
    id,
    user_id,
    token_hash,
    expires_at,
    consumed_at,
    created_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetPasswordResetTokenByHash :one
SELECT * FROM auth_password_reset_tokens
WHERE token_hash = $1
LIMIT 1;

-- name: ConsumePasswordResetToken :exec
UPDATE auth_password_reset_tokens
SET consumed_at = $2
WHERE id = $1;

-- name: InvalidatePasswordResetTokensForUser :exec
UPDATE auth_password_reset_tokens
SET consumed_at = $2
WHERE user_id = $1
  AND consumed_at IS NULL;
