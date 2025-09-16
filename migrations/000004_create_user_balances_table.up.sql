-- user_balances (materialized projection)
CREATE TABLE user_balances
(
    user_id    UUID PRIMARY KEY REFERENCES users (id),
    points     BIGINT      NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);