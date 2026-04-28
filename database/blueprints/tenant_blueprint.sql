-- ==========================================
-- DOMAIN: Profil Toko (Tenant Bridge)
-- ==========================================
CREATE TABLE store_profiles (
    store_id VARCHAR(255) PRIMARY KEY DEFAULT public.uuid_generate_v4(), 
    description TEXT,
    address TEXT,
    phone VARCHAR(20),
    email VARCHAR(100),
    tax_id VARCHAR(50),
    logo_url VARCHAR(255),
    receipt_footer TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- DOMAIN: Manajemen Pengguna Internal Toko
-- ==========================================

CREATE TABLE IF NOT EXISTS roles (
    id VARCHAR(50) PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    permissions JSONB 
);

CREATE TABLE IF NOT EXISTS employees (
    id VARCHAR(50) PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    user_id VARCHAR(50) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    gender public.gender,
    position VARCHAR(100),
    phone VARCHAR(20),
    email VARCHAR(100),
    address TEXT,
    photo VARCHAR(255),
    citizen_id VARCHAR(50),
    base_salary DECIMAL(15, 2),
    join_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    role_id VARCHAR(50) NOT NULL,
    FOREIGN KEY (role_id) REFERENCES roles(id)
);