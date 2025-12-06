CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed initial user for local development (password: password123)
INSERT INTO users (id, email, username, password, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'demo@example.com',
    'demo',
    '$2a$10$VsJaYdoUmPU2LBY.oLcZCeI.UuIkshR9OCdE3s9SXD5w8JVh2wQfa',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;

-- Bulk seed 100 local users (password: password123)
INSERT INTO users (id, email, username, password, created_at, updated_at)
SELECT
    uuid_generate_v4(),
    format('user%03s@example.com', lpad(gs.i::text, 3, '0')),
    format('user%03s', lpad(gs.i::text, 3, '0')),
    '$2a$10$VsJaYdoUmPU2LBY.oLcZCeI.UuIkshR9OCdE3s9SXD5w8JVh2wQfa',
    NOW(),
    NOW()
FROM generate_series(1, 100) AS gs(i)
ON CONFLICT (email) DO NOTHING;
