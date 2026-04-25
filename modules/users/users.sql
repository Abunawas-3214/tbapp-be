-- modules/users/users.sql

-- --- USER QUERIES ---

-- name: CreateUserByAdmin :one
INSERT INTO users (
    id, role_id, name, email, password_hash
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: CreateManyUsersByAdmin :copyfrom
INSERT INTO users (id, role_id, name, email, is_active)
VALUES ($1, $2, $3, $4, $5);

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users 
WHERE id = $1 LIMIT 1;

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

-- name: UpdateUsersStatus :exec
UPDATE users
SET is_active = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = ANY($1::varchar[]);

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: DeleteManyUsers :exec
DELETE FROM users
WHERE id = ANY($1::varchar[]);

-- --- ROLE QUERIES ---

-- name: CreateRole :one
INSERT INTO roles (id, name, permissions)
VALUES ($1, $2, $3) RETURNING *;

-- name: ListRoles :many
SELECT * FROM roles;

-- name: UpdateRole :one
UPDATE roles
SET name = $2, permissions = $3
WHERE id = $1 RETURNING *;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1;