-- name: CreatePermission :one
INSERT INTO access_permissions (
    id,
    key,
    description,
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

-- name: GetPermissionByKey :one
SELECT * FROM access_permissions
WHERE key = $1
LIMIT 1;

-- name: GetPermissionsByKeys :many
SELECT * FROM access_permissions
WHERE key = ANY($1::text[])
ORDER BY key;

-- name: ListPermissions :many
SELECT * FROM access_permissions
ORDER BY key;
