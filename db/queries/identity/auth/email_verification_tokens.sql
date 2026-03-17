-- name: CreateEmailVerificationToken :one
INSERT INTO auth_email_verification_tokens (
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

-- name: GetEmailVerificationTokenByHash :one
SELECT * FROM auth_email_verification_tokens
WHERE token_hash = $1
LIMIT 1;

-- name: ConsumeEmailVerificationToken :exec
UPDATE auth_email_verification_tokens
SET consumed_at = $2
WHERE id = $1;

-- name: InvalidateEmailVerificationTokensForUser :exec
UPDATE auth_email_verification_tokens
SET consumed_at = $2
WHERE user_id = $1
  AND consumed_at IS NULL;
