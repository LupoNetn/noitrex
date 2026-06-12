-- name: GetActiveWebhooksByOperatorID :one
SELECT * FROM webhook_endpoints WHERE operator_id = $1 AND active = true;


-- name: CreateDeliveryAttempt :one
INSERT INTO webhook_deliveries (endpoint_id,event_type,payload)
VALUES ($1,$2,$3)
RETURNING *;

-- name: UpdateDeliveryResult :one
UPDATE webhook_deliveries 
SET response_code = $1,
response_body = $2,
status = $3,
delivered_at = NOW(),
attempt_number = attempt_number + 1
WHERE id = $4
RETURNING *;

-- name: ListWebhookDeliveriesPaginated :many
SELECT d.id, d.endpoint_id, d.event_type, d.payload, d.response_code, d.response_body, d.attempt_number, d.delivered_at, d.created_at, d.status
FROM webhook_deliveries d
JOIN webhook_endpoints e ON d.endpoint_id = e.id
WHERE e.operator_id = $1
ORDER BY d.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateWebhookEndpointURL :one
UPDATE webhook_endpoints
SET url = $2
WHERE operator_id = $1 AND active = true
RETURNING *;

-- name: GetWebhookDeliveriesStats :one
SELECT 
    COUNT(d.id) as total_deliveries,
    COALESCE(SUM(CASE WHEN d.status = 'delivered' THEN 1 ELSE 0 END), 0)::bigint as successful_deliveries,
    COALESCE(SUM(CASE WHEN d.status = 'failed' THEN 1 ELSE 0 END), 0)::bigint as failed_deliveries
FROM webhook_deliveries d
JOIN webhook_endpoints e ON d.endpoint_id = e.id
WHERE e.operator_id = $1;
