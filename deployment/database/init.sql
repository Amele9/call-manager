CREATE TABLE IF NOT EXISTS calls
(
    id           SERIAL PRIMARY KEY,
    client_name  TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    description  TEXT NOT NULL,
    status       TEXT NOT NULL DEFAULT 'open',
    created_at   TIMESTAMPTZ   DEFAULT CURRENT_TIMESTAMP
)