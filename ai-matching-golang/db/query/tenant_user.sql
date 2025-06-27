-- name: AddUserToTenant :one
INSERT INTO tenant_users (
    tenant_id, user_id, role
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: RemoveUserFromTenant :exec
DELETE FROM tenant_users
WHERE tenant_id = $1 AND user_id = $2;

-- name: GetUsersByTenant :many
SELECT u.* FROM users u
INNER JOIN tenant_users tu ON u.id = tu.user_id
WHERE tu.tenant_id = $1
ORDER BY u.email;

-- name: GetTenantsByUser :many
SELECT t.* FROM tenants t
INNER JOIN tenant_users tu ON t.id = tu.tenant_id
WHERE tu.user_id = $1
ORDER BY t.name;

-- name: GetTenantUser :one
SELECT * FROM tenant_users
WHERE tenant_id = $1 AND user_id = $2
LIMIT 1;

-- name: UpdateUserRoleInTenant :one
UPDATE tenant_users
SET role = $3,
    updated_at = NOW()
WHERE tenant_id = $1 AND user_id = $2
RETURNING *;

-- name: ListTenantUsers :many
SELECT 
    tu.*,
    u.email,
    u.first_name,
    u.last_name
FROM tenant_users tu
INNER JOIN users u ON tu.user_id = u.id
WHERE tu.tenant_id = $1
ORDER BY u.email;