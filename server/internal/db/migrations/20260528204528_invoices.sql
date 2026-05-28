-- +goose Up
CREATE TYPE invoice_status AS ENUM (
    'draft',
    'issued',
    'paid',
    'void'
);

CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    operator_id UUID REFERENCES operators(id) ON DELETE CASCADE NOT NULL,
    customer_id UUID REFERENCES customers(id) ON DELETE CASCADE NOT NULL,
    amount_cents BIGINT NOT NULL,
    status invoice_status NOT NULL DEFAULT 'draft',
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    line_items JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(customer_id, period_start)
);

-- +goose Down
DROP TABLE IF EXISTS invoices;
DROP TYPE IF EXISTS invoice_status;
