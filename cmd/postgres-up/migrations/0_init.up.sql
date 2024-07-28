BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS outbox (
    id uuid DEFAULT uuid_generate_v4(),
    created_at timestamp with time zone DEFAULT now(),

    -- Message
    topic text NOT NULL,
    key bytea NOT NULL,
    value bytea NOT NULL,
    headers jsonb NOT NULL, -- e.g. [{"key": "Content-Type", "value_base64": "YXBwbGljYXRpb24vanNvbg=="}]

    PRIMARY KEY (id)
);

COMMIT;
