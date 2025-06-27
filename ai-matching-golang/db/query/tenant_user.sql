-- name: AddUserToTenant :one
INSERT INTO tenant_users (
    tenant_id, user_id, role
) VALUES (
    @tenant_id::uuid, @user_id::uuid, @role
)
RETURNING *;

-- name: RemoveUserFromTenant :exec
DELETE FROM tenant_users
WHERE tenant_id = @tenant_id::uuid AND user_id = @user_id::uuid;

-- name: GetUsersByTenant :many
SELECT u.* FROM users u
INNER JOIN tenant_users tu ON u.id = tu.user_id
WHERE tu.tenant_id = @tenant_id::uuid
ORDER BY u.email;

-- name: GetTenantsByUser :many
SELECT t.* FROM tenants t
INNER JOIN tenant_users tu ON t.id = tu.tenant_id
WHERE tu.user_id = @user_id::uuid
ORDER BY t.name;

-- name: GetTenantUser :one
SELECT * FROM tenant_users
WHERE tenant_id = @tenant_id::uuid AND user_id = @user_id::uuid
LIMIT 1;

-- name: UpdateUserRoleInTenant :one
UPDATE tenant_users
SET role = @role,
    updated_at = NOW()
WHERE tenant_id = @tenant_id::uuid AND user_id = @user_id::uuid
RETURNING *;

-- name: ListTenantUsers :many
SELECT 
    tu.*,
    u.email,
    u.first_name,
    u.last_name
FROM tenant_users tu
INNER JOIN users u ON tu.user_id = u.id
WHERE tu.tenant_id = @tenant_id::uuid
ORDER BY u.email;