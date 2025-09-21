-- name: ListRoles :many
SELECT * FROM roles ORDER BY name;

-- name: ListPermissions :many
SELECT * FROM permissions ORDER BY name;

-- name: GetPermissionsForRole :many
SELECT p.*
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = $1;

-- name: AddPermissionToRole :exec
INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: RemovePermissionFromRole :exec
DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2;