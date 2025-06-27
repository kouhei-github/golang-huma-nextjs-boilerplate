-- name: GetOrganization :one
SELECT * FROM organizations
WHERE id = @id::uuid LIMIT 1;

-- name: ListOrganizations :many
SELECT * FROM organizations
WHERE is_active = true
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: CreateOrganization :one
INSERT INTO organizations (
    name, description, is_active
) VALUES (
    @name, @description, @is_active
)
RETURNING *;

-- name: UpdateOrganization :one
UPDATE organizations
SET name = @name,
    description = @description,
    is_active = @is_active,
    updated_at = NOW()
WHERE id = @id::uuid
RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM organizations
WHERE id = @id::uuid;

-- name: GetOrganizationWithTenants :one
SELECT 
    o.*,
    COUNT(t.id) as tenant_count
FROM organizations o
LEFT JOIN tenants t ON o.id = t.organization_id AND t.is_active = true
WHERE o.id = @id::uuid
GROUP BY o.id;

-- name: GetTenantsByOrganization :many
SELECT * FROM tenants
WHERE organization_id = @organization_id::uuid AND is_active = true
ORDER BY name;

-- name: GetOrganizationByTenant :one
SELECT o.* FROM organizations o
INNER JOIN tenants t ON o.id = t.organization_id
WHERE t.id = @tenant_id::uuid
LIMIT 1;

-- name: CountOrganizations :one
SELECT COUNT(*) FROM organizations
WHERE is_active = true;