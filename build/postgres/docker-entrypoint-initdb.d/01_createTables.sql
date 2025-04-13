CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(64) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verified BOOLEAN DEFAULT false,
    phone VARCHAR(20) NOT NULL CHECK (
        phone ~ '^(\+7|7|8)[0-9]{10}$' AND  -- Основной формат
        phone !~ '.*[^0-9+].*' AND           -- Только цифры и +
        phone !~ '^8[0-9]{11}' AND           -- Запрет 12-значных номеров с 8
        phone !~ '^\+7[0-9]{11}'             -- Запрет 12-значных номеров с +7
    ),
    phone_verified BOOLEAN DEFAULT false,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    bio TEXT,
    telegram_username VARCHAR(64),
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_login_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT true,
    is_staff BOOLEAN DEFAULT false,
    
    CONSTRAINT valid_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$')
);


CREATE TABLE IF NOT EXISTS used_uuids (
    id SERIAL PRIMARY KEY,
    uuid_id UUID NOT NULL UNIQUE,
    used_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS user_companies (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    company_id UUID NOT NULL,
    isActive BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    logo_url TEXT,
    industry VARCHAR(100),
    employees INT DEFAULT 0,
    is_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    founded_date DATE,
    address TEXT,
    phone VARCHAR(20) NOT NULL CHECK (
        phone ~ '^(\+7|7|8)[0-9]{10}$' AND  -- Основной формат
        phone !~ '.*[^0-9+].*' AND           -- Только цифры и +
        phone !~ '^8[0-9]{11}' AND           -- Запрет 12-значных номеров с 8
        phone !~ '^\+7[0-9]{11}'             -- Запрет 12-значных номеров с +7
    ),
    email VARCHAR(255),
    tax_number VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS employee_company (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    company_id UUID,
    position_id UUID,
    is_active BOOLEAN DEFAULT true,
    is_online BOOLEAN DEFAULT false,
    last_activity_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID,
    lvl INTEGER,
    name TEXT,
    is_active BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID,
    name TEXT,
    description TEXT,
    avatar_url TEXT,
    parent_id UUID,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS employee_department (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL,
    department_id UUID NOT NULL,
    position_id UUID,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);
