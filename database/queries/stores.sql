-- name: CreateStore :one
INSERT INTO public.stores (name, slug, schema_name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateUserStoreAccess :exec
INSERT INTO public.user_store_access (user_id, store_id)
VALUES ($1, $2);

-- name: GetStoreBySlug :one
SELECT * FROM public.stores WHERE slug = $1 LIMIT 1;

-- name: UpdateStoreName :exec
UPDATE stores 
SET name = $2, updated_at = CURRENT_TIMESTAMP 
WHERE id = $1;