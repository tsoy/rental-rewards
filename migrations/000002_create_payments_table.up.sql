CREATE TABLE IF NOT EXISTS payments
(
    id           UUID PRIMARY KEY,
    user_id      UUID        NOT NULL REFERENCES users (id),
    amount_cents BIGINT      NOT NULL CHECK (amount_cents > 0),
    currency     TEXT        NOT NULL DEFAULT 'USD',
    status       TEXT        NOT NULL CHECK (status IN ('created', 'completed', 'failed')),
    external_ref TEXT, -- client idempotency key or external id
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ
);
-- enforce idempotency across the whole system (or choose per-user unique)
CREATE UNIQUE INDEX IF NOT EXISTS payments_external_ref_uidx ON
    payments (external_ref) WHERE external_ref IS NOT NULL