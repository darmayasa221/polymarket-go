package postgres

import "database/sql"

// schema defines all polymarket tables for PostgreSQL.
// All statements use CREATE TABLE/INDEX IF NOT EXISTS — safe to call repeatedly.
const schema = `
CREATE TABLE IF NOT EXISTS prices (
    id          TEXT        NOT NULL PRIMARY KEY,
    asset       TEXT        NOT NULL,
    source      TEXT        NOT NULL,
    value       TEXT        NOT NULL,
    rounded_at  TIMESTAMPTZ,
    received_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_prices_asset_received ON prices (asset, received_at DESC);

CREATE TABLE IF NOT EXISTS markets (
    id            TEXT        NOT NULL PRIMARY KEY,
    slug          TEXT        NOT NULL,
    asset         TEXT        NOT NULL,
    window_start  TIMESTAMPTZ NOT NULL,
    condition_id  TEXT        NOT NULL,
    up_token_id   TEXT        NOT NULL,
    down_token_id TEXT        NOT NULL,
    tick_size     TEXT        NOT NULL,
    fee_enabled   BOOLEAN     NOT NULL,
    active        BOOLEAN     NOT NULL,
    UNIQUE (asset, window_start)
);
CREATE INDEX IF NOT EXISTS idx_markets_active ON markets (active);

CREATE TABLE IF NOT EXISTS orders (
    id             TEXT        NOT NULL PRIMARY KEY,
    market_id      TEXT        NOT NULL,
    token_id       TEXT        NOT NULL,
    side           TEXT        NOT NULL,
    outcome        TEXT        NOT NULL,
    price          TEXT        NOT NULL,
    size           TEXT        NOT NULL,
    order_type     TEXT        NOT NULL,
    expiration     TIMESTAMPTZ,
    fee_rate_bps   BIGINT      NOT NULL,
    signature_type SMALLINT    NOT NULL,
    status         TEXT        NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_orders_market_status ON orders (market_id, status);

CREATE TABLE IF NOT EXISTS positions (
    id          TEXT        NOT NULL PRIMARY KEY,
    asset       TEXT        NOT NULL,
    token_id    TEXT        NOT NULL,
    outcome     TEXT        NOT NULL,
    size        TEXT        NOT NULL,
    avg_price   TEXT        NOT NULL,
    market_id   TEXT        NOT NULL,
    opened_at   TIMESTAMPTZ NOT NULL,
    closed_at   TIMESTAMPTZ,
    exit_price  TEXT
);
CREATE INDEX IF NOT EXISTS idx_positions_market ON positions (market_id);
CREATE INDEX IF NOT EXISTS idx_positions_closed  ON positions (closed_at);
`

// RunMigrations executes the schema DDL against db.
// Safe to call repeatedly — all statements are idempotent.
func RunMigrations(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}
