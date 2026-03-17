-- name: CreateSession :one
INSERT INTO auth_sessions (
    id,
    user_id,
    token_hash,
    user_agent,
    ip_address,
    expires_at,
    created_at
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

-- name: GetSessionByID :one
SELECT * FROM auth_sessions
WHERE id = $1
LIMIT 1;

-- name: GetSessionByTokenHash :one
SELECT * FROM auth_sessions
WHERE token_hash = $1
LIMIT 1;

-- name: RevokeSession :exec
UPDATE auth_sessions
SET revoked_at = $2
WHERE id = $1;

-- name: RevokeAllSessionsForUser :exec
UPDATE auth_sessions
SET revoked_at = $2
WHERE user_id = $1
  AND revoked_at IS NULL;

-- name: DeleteExpiredSessions :execrows
DELETE FROM auth_sessions
WHERE expires_at < $1;
