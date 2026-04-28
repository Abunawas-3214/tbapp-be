-- name: GetStoreProfile :one
SELECT 
    s.id as store_id,
    s.name as store_name,
    sp.description,
    sp.address,
    sp.phone,
    sp.email,
    sp.tax_id,
    sp.logo_url,
    sp.receipt_footer
FROM public.stores s
LEFT JOIN store_profiles sp ON s.id = sp.store_id
WHERE s.id = $1 LIMIT 1;

-- name: UpdateStoreProfile :one
UPDATE store_profiles
SET 
    description = $1,
    address = $2,
    phone = $3,
    email = $4,
    tax_id = $5,
    logo_url = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE store_id = $7
RETURNING *;