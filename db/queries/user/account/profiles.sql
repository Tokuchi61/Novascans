-- name: CreateProfile :one
INSERT INTO account_profiles (
    user_id,
    username,
    display_name,
    bio,
    avatar_path,
    banner_path,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetProfileByUserID :one
SELECT * FROM account_profiles
WHERE user_id = $1
LIMIT 1;

-- name: GetProfileByUsername :one
SELECT * FROM account_profiles
WHERE username = $1
LIMIT 1;

-- name: UpdateProfile :one
UPDATE account_profiles
SET username = $2,
    display_name = $3,
    bio = $4,
    avatar_path = $5,
    banner_path = $6,
    updated_at = $7
WHERE user_id = $1
RETURNING *;
