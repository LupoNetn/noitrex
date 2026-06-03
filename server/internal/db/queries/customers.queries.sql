-- name: CreateCustomer :one
INSERT INTO customers 
(operator_id, external_id, plan_name, name, email)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCustomerByID :one
SELECT * FROM customers WHERE id = $1;

-- name: GetCustomerByExternalID :one
SELECT * FROM customers WHERE operator_id = $1 AND external_id = $2;

-- name: GetCustomerByEmail :one
SELECT * FROM customers WHERE operator_id = $1 AND email = $2;

-- name: UpdateCustomerPlan :one
UPDATE customers SET plan_name = $1, updated_at = NOW() WHERE operator_id = $2 AND id = $3 RETURNING *;

-- name: ListCustomers :many
SELECT * FROM customers WHERE operator_id = $1 ORDER BY created_at DESC;