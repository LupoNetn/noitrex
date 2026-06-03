-- name: CreatePlan :one
INSERT INTO plans 
(operator_id,name,pricing_model, unit_price_cent,tiers, billing_period) 
VALUES ($1,$2,$3,$4,$5,$6) 
RETURNING *;

-- name: GetPlan :one
SELECT * FROM plans WHERE id = $1;

-- name: ListPlans :many
SELECT * FROM plans WHERE operator_id = $1;

-- name: GetPlanByName :one
SELECT * FROM plans WHERE operator_id = $1 AND name = $2;

-- name: UpdatePlan :one
UPDATE plans SET 
    name = COALESCE($1, name),
    pricing_model = COALESCE($2, pricing_model),
    unit_price_cent = COALESCE($3, unit_price_cent),
    tiers = COALESCE($4, tiers),
    billing_period = COALESCE($5, billing_period),
    updated_at = NOW()
WHERE id = $6 RETURNING *;

-- name: DeletePlan :exec
DELETE FROM plans WHERE id = $1;