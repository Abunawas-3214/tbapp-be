-- name: GetStoreProfile :one
SELECT * FROM store_profiles LIMIT 1;

-- name: UpdateStoreProfile :exec
UPDATE store_profiles 
SET name = $2, description = $3, address = $4, phone = $5, email = $6, tax_id = $7, logo_url = $8, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;