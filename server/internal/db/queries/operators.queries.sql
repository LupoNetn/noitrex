-- name: CreateOperator :one
INSERT INTO operators (
    name, 
    description, 
    logo_url, 
    api_key_hash, 
    webhook_secret, 
    support_email, 
    website_url, 
    status, 
    password_hash
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetOperator :one
SELECT * FROM operators WHERE id = $1;

-- name: GetOperatorByEmail :one
SELECT * FROM operators WHERE support_email = $1;

-- name: ListOperators :many
SELECT * FROM operators ORDER BY created_at DESC;

-- name: UpdateOperator :one
UPDATE operators
SET 
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    logo_url = COALESCE($4, logo_url),
    support_email = COALESCE($5, support_email),
    website_url = COALESCE($6, website_url),
    password_hash = COALESCE($7, password_hash),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateOperatorStatus :exec
UPDATE operators
SET status = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteOperator :exec
DELETE FROM operators WHERE id = $1;