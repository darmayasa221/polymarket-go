CREATE TABLE IF NOT EXISTS users (
    id              TEXT        NOT NULL PRIMARY KEY,
    username        TEXT        NOT NULL UNIQUE,
    email           TEXT        NOT NULL UNIQUE,
    hashed_password TEXT        NOT NULL,
    full_name       TEXT        NOT NULL,
    created_at      DATETIME    NOT NULL,
    updated_at      DATETIME    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email    ON users(email);
