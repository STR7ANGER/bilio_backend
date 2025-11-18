BEGIN;

CREATE TABLE IF NOT EXISTS promocodes (
    code TEXT PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    used_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_promocodes_used_at ON promocodes(used_at) WHERE used_at IS NULL;

COMMIT;

