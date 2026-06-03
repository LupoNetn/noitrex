-- +goose Up
ALTER TABLE plans
ADD CONSTRAINT unique_name_per_operator UNIQUE (operator_id, name);

-- +goose Down
ALTER TABLE plans
DROP CONSTRAINT unique_name_per_operator;
