-- +goose Up
CREATE TABLE IF NOT EXISTS usage_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID REFERENCES customers(id),
    operator_id UUID REFERENCES operators(id),
    event_name TEXT NOT NULL,
    quantity BIGINT NOT NULL,
    idempotency_key TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(customer_id, idempotency_key)
);

-- +goose Down
DROP TABLE IF EXISTS usage_events;

