-- +goose Up
ALTER TABLE customers DROP COLUMN plan_id;
ALTER TABLE plans DROP CONSTRAINT unique_name_per_operator;
ALTER TABLE plans DROP COLUMN name;
ALTER TABLE plans ADD COLUMN name TEXT UNIQUE NOT NULL;
ALTER TABLE plans DROP COLUMN operator_id;
ALTER TABLE plans ADD COLUMN operator_id UUID NOT NULL REFERENCES operators(id);
ALTER TABLE customers ADD COLUMN plan_name TEXT NOT NULL REFERENCES plans(name);

-- +goose Down
ALTER TABLE customers DROP COLUMN plan_name;
