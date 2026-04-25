-- modules/users/users.sql

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users 
WHERE id = $1 LIMIT 1;

-- name: CreateUserByAdmin :one
INSERT INTO users (
    id, role_id, name, email
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: CreateRole :one
INSERT INTO roles (
    id, name, permissions
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: ListUsers :many
SELECT * FROM users 
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET 
    name = $2,
    role_id = $3,
    is_active = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;