-- name: GetTenant :one
SELECT * FROM tenants
WHERE id = $1 LIMIT 1;

-- name: GetTenantBySubdomain :one
SELECT * FROM tenants
WHERE subdomain = $1 AND is_active = true
LIMIT 1;

-- name: ListTenantsByOrganization :many
SELECT * FROM tenants
WHERE organization_id = $1 AND is_active = true
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: CreateTenant :one
INSERT INTO tenants (
    organization_id, name, subdomain, is_active
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateTenant :one
UPDATE tenants
SET name = $2,
    subdomain = $3,
    is_active = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteTenant :exec
DELETE FROM tenants
WHERE id = $1;

-- name: GetTenantWithUserCount :one
SELECT 
    t.*,
    COUNT(DISTINCT tu.user_id) as user_count
FROM tenants t
LEFT JOIN tenant_users tu ON t.id = tu.tenant_id
WHERE t.id = $1
GROUP BY t.id;

-- name: GetTenantsByUserID :many
SELECT t.*, tu.role
FROM tenants t
INNER JOIN tenant_users tu ON t.id = tu.tenant_id
WHERE tu.user_id = $1 AND t.is_active = true
ORDER BY t.name;

-- name: CheckUserBelongsToTenant :one
SELECT EXISTS(
    SELECT 1 FROM tenant_users
    WHERE tenant_id = $1 AND user_id = $2
) as belongs;