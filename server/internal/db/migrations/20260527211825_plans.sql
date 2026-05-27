-- +goose Up
CREATE TYPE pricing_type AS ENUM (
    'flat',
    'tiered',
    'per_unit'
);

CREATE TYPE billing_period_type AS ENUM (
    'monthly',
    'yearly'
);

CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    operator_id UUID NOT NULL REFERENCES operators(id),
    name TEXT NOT NULL,
    pricing_model pricing_type NOT NULL,
    unit_price_cent BIGINT,
    tiers JSONB,
    billing_period billing_period_type NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_pricing_model CHECK (
        (pricing_model = 'tiered' AND unit_price_cent IS NULL AND tiers IS NOT NULL) OR 
        (pricing_model != 'tiered' AND unit_price_cent IS NOT NULL AND tiers IS NULL)
    )
);

-- +goose Down
DROP TABLE IF EXISTS plans;
DROP TYPE IF EXISTS pricing_type;
DROP TYPE IF EXISTS billing_period;