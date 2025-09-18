-- name: CreateRole :one
INSERT INTO roles (name) VALUES ($1) RETURNING *;

-- name: CreatePermission :one
INSERT INTO permissions (name) VALUES ($1) RETURNING *;

-- name: GetUserPermissions :many
-- Retrieves all permission names for a given user ID by joining through
-- user_roles and role_permissions.
SELECT p.name
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
JOIN user_roles ur ON rp.role_id = ur.role_id
WHERE ur.user_id = $1;