-- Создание таблицы users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY  UNIQUE,
    login VARCHAR(63) NOT NULL UNIQUE,
    password VARCHAR(63) NOT NULL,
    name TEXT,
    surname TEXT,
    bio TEXT,
    phone VARCHAR(11),
    telegram VARCHAR(63),
    mail VARCHAR(63) NOT NULL UNIQUE
);

-- Создание таблицы companies
CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY UNIQUE,
    owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name  TEXT,
    description TEXT
);

-- Создание таблицы company_info
CREATE TABLE IF NOT EXISTS company_info (
    id UUID PRIMARY KEY,
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE
);

-- Создание таблицы positions
CREATE TABLE IF NOT EXISTS positions (
    id UUID PRIMARY KEY,
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    lvl INTEGER NOT NULL, -- 0 - owner -- 1 amdin -- 2 employee -- no role
    name TEXT NOT NULL
);

-- Создание таблицы employee_company
CREATE TABLE IF NOT EXISTS employee_company (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    position UUID REFERENCES positions(id) ON DELETE CASCADE
);

-- Создание таблицы departments
CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY,
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    owner_employee_id UUID REFERENCES employee_company(id) ON DELETE CASCADE
);

-- Создание таблицы employee_department
CREATE TABLE IF NOT EXISTS employee_department (
    id UUID PRIMARY KEY,
    employee_id UUID REFERENCES employee_company(id) ON DELETE CASCADE,
    department_id UUID REFERENCES departments(id) ON DELETE CASCADE
);

-- Создание таблицы invites
CREATE TABLE IF NOT EXISTS invites (
    id UUID PRIMARY KEY,
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    mail VARCHAR(63),
    link TEXT,
    timeout DATE
);

-- Создание таблицы access_options
CREATE TABLE IF NOT EXISTS access_options (
    id UUID PRIMARY KEY,
    option INTEGER NOT NULL UNIQUE
);

-- Создание таблицы access
CREATE TABLE IF NOT EXISTS access (
    id SERIAL PRIMARY KEY,
    employee_id UUID REFERENCES employee_department(id) ON DELETE CASCADE,
    access_option_id UUID REFERENCES access_options(id) ON DELETE CASCADE
);

-- Добавление индексов для улучшения производительности
CREATE INDEX idx_users_login ON users(login);
CREATE INDEX idx_users_mail ON users(mail);
CREATE INDEX idx_companies_owner ON companies(owner_id);
CREATE INDEX idx_employee_company_user ON employee_company(user_id);
CREATE INDEX idx_employee_company_company ON employee_company(company_id);
CREATE INDEX idx_employee_department_employee ON employee_department(employee_id);
CREATE INDEX idx_employee_department_department ON employee_department(department_id);
CREATE INDEX idx_departments_company ON departments(company_id);
CREATE INDEX idx_positions_company ON positions(company_id);
CREATE INDEX idx_invites_company ON invites(company_id);
CREATE INDEX idx_access_employee ON access(employee_id);
CREATE INDEX idx_access_option ON access(access_option_id);