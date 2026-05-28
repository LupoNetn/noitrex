-- name: CreateOperator :one
INSERT INTO operators 
(name, description, status, api_key_hash, webhook_secret, support_email, website_url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;


-- name: GetOperator :one
SELECT * FROM operators WHERE id = $1;