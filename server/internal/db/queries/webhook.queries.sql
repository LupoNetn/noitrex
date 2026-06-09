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



