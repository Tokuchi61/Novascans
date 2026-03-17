-- name: AssignSubRoleToUser :exec
INSERT INTO access_user_sub_roles (
    user_id,
    sub_role_id,
    assigned_at
) VALUES (
    $1,
    $2,
    $3
)
ON CONFLICT (user_id, sub_role_id) DO NOTHING;

-- name: RemoveSubRoleFromUser :exec
DELETE FROM access_user_sub_roles
WHERE user_id = $1
  AND sub_role_id = $2;

-- name: ListSubRolesForUser :many
SELECT
    sr.id,
    sr.key,
    sr.name,
    sr.description,
    sr.created_at,
    sr.updated_at
FROM access_user_sub_roles usr
INNER JOIN access_sub_roles sr ON sr.id = usr.sub_role_id
WHERE usr.user_id = $1
ORDER BY sr.key;

-- name: ListPermissionKeysForUser :many
SELECT DISTINCT p.key
FROM access_user_sub_roles usr
INNER JOIN access_sub_role_permissions srp ON srp.sub_role_id = usr.sub_role_id
INNER JOIN access_permissions p ON p.id = srp.permission_id
WHERE usr.user_id = $1
ORDER BY p.key;
