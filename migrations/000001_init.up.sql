CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_role AS ENUM ('admin', 'analyst', 'viewer');
CREATE TYPE record_type AS ENUM ('income', 'expense');

CREATE TABLE users (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT        NOT NULL,
    email      TEXT        NOT NULL UNIQUE,
    password   TEXT        NOT NULL,
    role       user_role   NOT NULL DEFAULT 'viewer',
    is_active  BOOLEAN     NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE financial_records (
    id          UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount      NUMERIC(12,2) NOT NULL CHECK (amount > 0),
    type        record_type   NOT NULL,
    category    TEXT          NOT NULL,
    description TEXT,
    date        DATE          NOT NULL,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_records_user_id  ON financial_records(user_id);
CREATE INDEX idx_records_date     ON financial_records(date);
CREATE INDEX idx_records_type     ON financial_records(type);
CREATE INDEX idx_records_category ON financial_records(category);