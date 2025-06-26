-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
    email, password_hash, first_name, last_name, organization_id, tenant_id
) VALUES (
    $1, $2, $3, $4, $5, $6
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