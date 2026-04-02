CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auth_user_id TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    role TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_role_check CHECK (role IN ('viewer', 'analyst', 'admin')),
    CONSTRAINT users_status_check CHECK (status IN ('active', 'inactive'))
);

CREATE INDEX idx_users_auth_user_id ON users (auth_user_id);
CREATE INDEX idx_users_role ON users (role);
CREATE INDEX idx_users_status ON users (status);

---- create above / drop below ----

DROP TABLE IF EXISTS users;
