-- name: GetUser :one
SELECT * FROM users
WHERE id = @id::uuid LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = @email LIMIT 1;

-- name: GetUserByCognitoID :one
SELECT * FROM users
WHERE cognito_id = @cognito_id LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
    cognito_id, email, first_name, last_name
) VALUES (
             @cognito_id, @email, @first_name, @last_name
         )
    RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET email = @email,
    first_name = @first_name,
    last_name = @last_name,
    updated_at = NOW()
WHERE id = @id::uuid
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = @id::uuid;

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
WHERE u.id = @id::uuid
GROUP BY u.id;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: CountUsersNotInTenant :one
SELECT COUNT(*)
FROM users u
WHERE u.id NOT IN (
    SELECT user_id FROM tenant_users WHERE tenant_id = $1
);