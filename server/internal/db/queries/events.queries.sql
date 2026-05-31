-- name: CreateUsageEvent :one
INSERT INTO usage_events (customer_id, operator_id, event_name, quantity, idempotency_key)
VALUES ($1, $2, $3, $4, $5) RETURNING *;


-- name: GetUsageEventByIdempotencyKey :one
SELECT * FROM usage_events WHERE idempotency_key = $1;

-- name: GetUsageEventByID :one
SELECT * FROM usage_events WHERE id = $1;

-- name: GetUsageEventByCustomerID :one
SELECT * FROM usage_events WHERE customer_id = $1;

-- name: ListUsageEventsByCustomerID :many
SELECT * FROM usage_events WHERE customer_id = $1;

-- name: ListUsageEventsByEventName :many
SELECT * FROM usage_events WHERE event_name = $1;

-- name: GetUsageEventByCustomerIDAndEventName :one
SELECT * FROM usage_events WHERE customer_id = $1 AND event_name = $2;