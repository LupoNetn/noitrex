-- +goose Up
ALTER TABLE webhook_deliveries ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'pending';

-- +goose Down
ALTER TABLE webhook_deliveries DROP COLUMN status;
