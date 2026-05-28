-- +goose Up
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    operator_id UUID NOT NULL REFERENCES operators(id),
    external_id UUID NOT NULL,
    plan_id UUID NOT NULL REFERENCES plans(id),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    period_start TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS customers;
