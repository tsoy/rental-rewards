CREATE TABLE IF NOT EXISTS users
(
    id         UUID PRIMARY KEY,
    email      TEXT UNIQUE NOT NULL,
    full_name  TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);