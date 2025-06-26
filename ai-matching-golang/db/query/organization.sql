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