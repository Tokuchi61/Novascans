-- name: CreateSubRole :one
INSERT INTO access_sub_roles (
    id,
    key,
    name,
    description,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetSubRoleByID :one
SELECT * FROM access_sub_roles
WHERE id = $1
LIMIT 1;

-- name: GetSubRoleByKey :one
SELECT * FROM access_sub_roles
WHERE key = $1
LIMIT 1;

-- name: ListSubRoles :many
SELECT * FROM access_sub_roles
ORDER BY key;

-- name: AttachPermissionToSubRole :exec
INSERT INTO access_sub_role_permissions (
    sub_role_id,
    permission_id,
    created_at
) VALUES (
    $1,
    $2,
    $3
)
ON CONFLICT (sub_role_id, permission_id) DO NOTHING;

-- name: ListSubRolePermissionLinks :many
SELECT
    sr.id AS sub_role_id,
    sr.key AS sub_role_key,
    sr.name AS sub_role_name,
    sr.description AS sub_role_description,
    sr.created_at AS sub_role_created_at,
    sr.updated_at AS sub_role_updated_at,
    p.id AS permission_id,
    p.key AS permission_key,
    p.description AS permission_description,
    p.created_at AS permission_created_at,
    p.updated_at AS permission_updated_at
FROM access_sub_roles sr
LEFT JOIN access_sub_role_permissions srp ON srp.sub_role_id = sr.id
LEFT JOIN access_permissions p ON p.id = srp.permission_id
ORDER BY sr.key, p.key;
