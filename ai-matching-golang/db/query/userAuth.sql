-- name: CreateUserAuth :one
INSERT INTO user_auth (
    user_id, refresh_token, expires_at
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetUserAuthByToken :one
SELECT * FROM user_auth
WHERE refresh_token = $1 AND expires_at > NOW()
LIMIT 1;

-- name: GetUserAuthByUserID :one
SELECT * FROM user_auth
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateUserAuth :one
UPDATE user_auth
SET refresh_token = $2,
    expires_at = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUserAuth :exec
DELETE FROM user_auth
WHERE id = $1;

-- name: DeleteExpiredUserAuth :exec
DELETE FROM user_auth
WHERE expires_at < NOW();