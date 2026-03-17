-- name: CreateSettings :one
INSERT INTO account_settings (
    user_id,
    locale,
    timezone,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetSettingsByUserID :one
SELECT * FROM account_settings
WHERE user_id = $1
LIMIT 1;

-- name: UpdateSettings :one
UPDATE account_settings
SET locale = $2,
    timezone = $3,
    updated_at = $4
WHERE user_id = $1
RETURNING *;
