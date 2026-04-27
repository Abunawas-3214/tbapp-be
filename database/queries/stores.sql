-- name: CreateStore :one
INSERT INTO public.stores (name, slug, schema_name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateUserStoreAccess :exec
INSERT INTO public.user_store_access (user_id, store_id)
VALUES ($1, $2);

-- name: GetStoreBySlug :one
SELECT * FROM public.stores WHERE slug = $1 LIMIT 1;