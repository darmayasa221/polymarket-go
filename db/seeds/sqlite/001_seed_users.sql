-- Seed data for development and testing.
-- Passwords are bcrypt hashes of "password123" (cost 12).

INSERT OR IGNORE INTO users (id, username, email, hashed_password, full_name, created_at, updated_at)
VALUES (
    'usr_01h2xkm8p0000000000000000',
    'admin',
    'admin@example.com',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj4J/HS.iK8i',
    'Admin User',
    '2026-01-01 00:00:00',
    '2026-01-01 00:00:00'
);
