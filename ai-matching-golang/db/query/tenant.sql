-- name: GetTenant :one
SELECT * FROM tenants
WHERE id = @id::uuid LIMIT 1;

-- name: GetTenantBySubdomain :one
SELECT * FROM tenants
WHERE subdomain = @subdomain AND is_active = true
LIMIT 1;

-- name: ListTenantsByOrganization :many
SELECT * FROM tenants
WHERE organization_id = @organization_id::uuid AND is_active = true
ORDER BY id
    LIMIT $1 OFFSET $2;

-- name: CreateTenant :one
INSERT INTO tenants (
    organization_id, name, subdomain, is_active
) VALUES (
    @organization_id::uuid, @name, @subdomain, @is_active
)
RETURNING *;

-- name: UpdateTenant :one
UPDATE tenants
SET name = @name,
    subdomain = @subdomain,
    is_active = @is_active,
    updated_at = NOW()
WHERE id = @id::uuid
RETURNING *;

-- name: DeleteTenant :exec
DELETE FROM tenants
WHERE id = @id::uuid;

-- name: GetTenantWithUserCount :one
SELECT 
    t.*,
    COUNT(DISTINCT tu.user_id) as user_count
FROM tenants t
LEFT JOIN tenant_users tu ON t.id = tu.tenant_id
WHERE t.id = @id::uuid
GROUP BY t.id;

-- name: GetTenantsByUserID :many
SELECT t.*, tu.role
FROM tenants t
INNER JOIN tenant_users tu ON t.id = tu.tenant_id
WHERE tu.user_id = @user_id::uuid AND t.is_active = true
ORDER BY t.name;

-- name: CheckUserBelongsToTenant :one
SELECT EXISTS(
    SELECT 1 FROM tenant_users
    WHERE tenant_id = @tenant_id::uuid AND user_id = @user_id::uuid
) as belongs;