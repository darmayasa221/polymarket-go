CREATE TABLE IF NOT EXISTS tokens (
    id          TEXT        NOT NULL PRIMARY KEY,
    user_id     TEXT        NOT NULL,
    value       TEXT        NOT NULL UNIQUE,
    type        TEXT        NOT NULL,
    purpose     TEXT        NOT NULL,
    expires_at  DATETIME    NOT NULL,
    created_at  DATETIME    NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_tokens_value   ON tokens(value);
CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id);
