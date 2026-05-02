-- name: GetTenantEmployeeRole :one
-- Mengambil data lengkap karyawan beserta role dan permission-nya
-- Digunakan saat proses Select Store untuk mengisi klaim JWT Tenant
SELECT 
    e.user_id,
    e.is_active,
    r.id AS role_id,
    r.name AS role_name,
    r.permissions
FROM employees e
JOIN roles r ON e.roleId = r.id
WHERE e.user_id = $1 AND e.isActive = true 
LIMIT 1;

-- name: ListEmployeesWithRoles :many
-- Mengambil daftar semua karyawan di dalam toko/tenant
SELECT 
    e.id, 
    e.full_name, 
    e.position, 
    e.is_active, 
    r.name as role_name
FROM employees e
LEFT JOIN roles r ON e.roleId = r.id
ORDER BY e.full_name ASC;