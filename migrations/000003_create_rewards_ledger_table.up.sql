CREATE TABLE rewards_ledger
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    UUID        NOT NULL REFERENCES users (id),
    payment_id UUID REFERENCES payments (id),
    points     BIGINT      NOT NULL,
    reason     TEXT        NOT NULL, -- e.g. PAYMENT_COMPLETED
    event_id   TEXT        NOT NULL, -- Pub/Sub message id for idempotency
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (event_id)
);