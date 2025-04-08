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
