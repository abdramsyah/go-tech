CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    nik VARCHAR(16) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    legal_name VARCHAR(100) NOT NULL,
    birth_place VARCHAR(50) NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    birth_date VARCHAR(50) NOT NULL,
    salary NUMERIC(15, 2) NOT NULL,
    ktp_photo TEXT NOT NULL,
    selfie_photo TEXT NOT NULL,
    status VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    role_id INTEGER NOT NULL,
    password_hash TEXT NOT NULL,
    created_by BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_by BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Ensure email uniqueness
CREATE UNIQUE INDEX idx_users_email ON users(email);

-- Ensure nik uniqueness
CREATE UNIQUE INDEX idx_users_nik ON users(nik);
