-- name: CreateRole :one
INSERT INTO roles (id, name, permissions)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateEmployee :one
INSERT INTO employees (id, user_id, full_name, role_id)
VALUES ($1, $2, $3, $4)
RETURNING *;