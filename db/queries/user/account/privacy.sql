-- name: CreatePrivacySettings :one
INSERT INTO account_privacy_settings (
    user_id,
    profile_visibility,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetPrivacySettingsByUserID :one
SELECT * FROM account_privacy_settings
WHERE user_id = $1
LIMIT 1;

-- name: UpdatePrivacySettings :one
UPDATE account_privacy_settings
SET profile_visibility = $2,
    updated_at = $3
WHERE user_id = $1
RETURNING *;
