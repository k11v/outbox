BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS outbox (
    id uuid DEFAULT uuid_generate_v4(),
    topic text NOT NULL,
    key bytea NOT NULL,
    value bytea NOT NULL,
    headers jsonb NOT NULL, -- e.g. [{"key": "Content-Type", "value": "YXBwbGljYXRpb24vanNvbg=="}]
    PRIMARY KEY (id)
);

COMMIT;
