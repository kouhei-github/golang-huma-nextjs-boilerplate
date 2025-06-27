-- name: GetOrganization :one
SELECT * FROM organizations
WHERE id = $1 LIMIT 1;

-- name: ListOrganizations :many
SELECT * FROM organizations
WHERE is_active = true
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: CreateOrganization :one
INSERT INTO organizations (
    name, description, is_active
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateOrganization :one
UPDATE organizations
SET name = $2,
    description = $3,
    is_active = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM organizations
WHERE id = $1;

-- name: GetOrganizationWithTenants :one
SELECT 
    o.*,
    COUNT(t.id) as tenant_count
FROM organizations o
LEFT JOIN tenants t ON o.id = t.organization_id AND t.is_active = true
WHERE o.id = $1
GROUP BY o.id;

-- name: GetTenantsByOrganization :many
SELECT * FROM tenants
WHERE organization_id = $1 AND is_active = true
ORDER BY name;

-- name: GetOrganizationByTenant :one
SELECT o.* FROM organizations o
INNER JOIN tenants t ON o.id = t.organization_id
WHERE t.id = $1
LIMIT 1;