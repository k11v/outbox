BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS message_infos (
    id uuid DEFAULT uuid_generate_v4(),
    value_length integer NOT NULL
);

CREATE TABLE IF NOT EXISTS outbox_messages (
    id uuid DEFAULT uuid_generate_v4(),
    created_at timestamp with time zone DEFAULT now(),
    status text NOT NULL,

    -- Message.
    -- For simplicity, we're assuming that all keys and values are strings.
    -- In a real-world scenario, you might want to use bytea or base64.
    topic text NOT NULL,
    key text NOT NULL,
    value text NOT NULL,
    headers jsonb NOT NULL, -- e.g. [{"key": "Content-Type", "value": "application/json"}]

    PRIMARY KEY (id),
    CHECK (status IN ('undelivered', 'delivered'))
);

COMMIT;
