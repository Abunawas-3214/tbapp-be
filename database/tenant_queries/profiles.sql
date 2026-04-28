-- name: GetStoreProfile :one
SELECT * FROM store_profiles LIMIT 1;

-- name: UpdateStoreProfile :one
UPDATE store_profiles
SET 
    name = $1,
    description = $2,
    address = $3,
    phone = $4,
    email = $5,
    tax_id = $6,
    logo_url = $7,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $8
RETURNING *;