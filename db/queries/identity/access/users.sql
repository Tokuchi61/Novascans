-- name: UpdateUserBaseRole :exec
UPDATE users
SET base_role = $2,
    updated_at = $3
WHERE id = $1;
