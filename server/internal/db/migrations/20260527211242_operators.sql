-- +goose Up
CREATE TYPE operating_status AS ENUM (
    'active',
    'inactive',
    'maintenance',
    'deprecated'
);

CREATE TABLE IF NOT EXISTS operators (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    logo_url TEXT,
    api_key_hash TEXT NOT NULL,
    webhook_secret TEXT NOT NULL,
    support_email TEXT,
    website_url TEXT,
    status operating_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- +goose Down
DROP TABLE IF EXISTS operators;
DROP TYPE operating_status;
