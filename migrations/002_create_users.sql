CREATE TABLE IF NOT EXISTS users (
    id          UUID PRIMARY KEY,
    email       TEXT UNIQUE NOT NULL,
    pass_hash   TEXT NOT NULL,
    role        TEXT NOT NULL CHECK (role IN ('employee','moderator')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
