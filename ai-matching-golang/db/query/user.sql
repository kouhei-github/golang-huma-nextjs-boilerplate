-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByCognitoID :one
SELECT * FROM users
WHERE cognito_id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
    cognito_id, email, first_name, last_name
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET email = $2,
    first_name = $3,
    last_name = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUsersNotInTenant :many
SELECT u.*
FROM users u
WHERE u.id NOT IN (
    SELECT user_id FROM tenant_users WHERE tenant_id = $1
)
ORDER BY u.email
LIMIT $2 OFFSET $3;

-- name: GetUserWithTenants :one
SELECT 
    u.*,
    COUNT(DISTINCT tu.tenant_id) as tenant_count
FROM users u
LEFT JOIN tenant_users tu ON u.id = tu.user_id
WHERE u.id = $1
GROUP BY u.id;