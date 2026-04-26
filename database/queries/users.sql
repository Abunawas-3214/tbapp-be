-- ==========================================
-- DOMAIN: Identitas Pengguna (Global)
-- ==========================================

-- name: CreateBaseUser :one
-- Digunakan oleh: Modul Sysadmin (Global) & Modul Tenant (Onboarding)
INSERT INTO public.users (
    name, 
    email, 
    password_hash, 
    image
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUserByEmail :one
-- Digunakan untuk: Proses login awal
SELECT * FROM public.users
WHERE email = $1 LIMIT 1;

-- name: GetUserById :one
-- Digunakan untuk: Fetch data profil user
SELECT * FROM public.users
WHERE id = $1 LIMIT 1;

-- name: UpdateUserStatus :exec
-- Digunakan untuk: Memblokir atau mengaktifkan user secara global
UPDATE public.users
SET is_active = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateUserProfile :one
-- Digunakan oleh: Modul Shared (Profile)
UPDATE public.users
SET 
    name = $2,
    image = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
-- Catatan: Akan memicu CASCADE ke system_admins & user_store_access
DELETE FROM public.users
WHERE id = $1;


-- ==========================================
-- DOMAIN: Otoritas Admin Sistem (Global)
-- ==========================================

-- name: CreateSystemAdmin :one
-- Digunakan oleh: Modul Sysadmin
INSERT INTO public.system_admins (
    user_id, 
    level
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetAdminAccess :one
-- Digunakan oleh: Middleware untuk cek apakah user punya akses Superadmin
SELECT level FROM public.system_admins
WHERE user_id = $1 LIMIT 1;

-- name: GetSystemAdminByID :one
SELECT u.id, u.name, u.email, u.is_active, sa.level, u.created_at
FROM public.users u
JOIN public.system_admins sa ON u.id = sa.user_id
WHERE u.id = $1 LIMIT 1;

-- name: ListAllSystemAdmins :many
-- Menampilkan daftar admin beserta info user-nya (Join)
SELECT 
    u.id, u.name, u.email, u.image, u.is_active,
    sa.level, sa.created_at as admin_since
FROM public.users u
JOIN public.system_admins sa ON u.id = sa.user_id
ORDER BY sa.created_at DESC;

-- name: UpdateSystemAdmin :exec
UPDATE public.users 
SET name = $2, is_active = $3
WHERE id = $1;

-- name: UpdateAdminLevel :exec
UPDATE public.system_admins
SET level = $2
WHERE user_id = $1;

-- name: DeleteSystemAdmin :exec
DELETE FROM public.users WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE public.users SET password_hash = $2 WHERE id = $1;

-- name: UpdateUserEmail :exec
UPDATE public.users 
SET email = $2 
WHERE id = $1;


-- ==========================================
-- DOMAIN: Otoritas Akses Toko (Tenant Bridge)
-- ==========================================

-- name: AssignUserToStore :one
-- Digunakan oleh: Modul Tenant (saat buat toko) & Modul Employee (saat tambah staf)
INSERT INTO public.user_store_access (
    user_id, 
    store_id, 
    is_active
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUserStores :many
-- Digunakan untuk: Layar "Pilih Toko" setelah login
SELECT 
    s.id, s.name, s.slug, s.schema_name, 
    usa.is_active as access_status
FROM public.stores s
JOIN public.user_store_access usa ON s.id = usa.store_id
WHERE usa.user_id = $1 AND s.is_active = true;