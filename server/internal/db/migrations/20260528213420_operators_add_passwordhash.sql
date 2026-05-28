-- +goose Up
ALTER TABLE operators 
ADD COLUMN password_hash TEXT;

-- +goose Down
ALTER TABLE operators 
DROP COLUMN password_hash;
