-- name: GetCustomersDueForBillingWithoutInvoice :many
SELECT 
c.id AS customer_id,
c.operator_id,
c.name AS customer_name,
c.period_start,
p.name AS plan_name,
p.pricing_model AS plan_pricing_model,
p.tiers AS plan_tiers,
p.unit_price_cent AS plan_unit_price,
p.billing_period AS plan_billing_period,
ua.total_units,
ua.event_name
FROM customers c
JOIN plans p ON p.name = c.plan_name
LEFT JOIN usage_aggregates ua ON ua.customer_id = c.id AND ua.period_start = c.period_start
WHERE (
    (p.billing_period = 'monthly' AND c.period_start + INTERVAL '1 month' <= NOW()) 
    OR
    (p.billing_period = 'yearly' AND c.period_start + INTERVAL '1 year' <= NOW())
)
AND NOT EXISTS (
    SELECT 1
    FROM invoices i
    WHERE i.customer_id = c.id
    AND i.period_start = c.period_start
    AND i.status != 'void'
)
ORDER BY c.id;

-- name: CreateNewInvoice :one
INSERT INTO invoices (
    operator_id,
    customer_id,
    amount_cents,
    status,
    period_start,
    period_end,
    line_items
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetInvoiceByID :one
SELECT * FROM invoices 
WHERE id = $1;

-- name: ListInvoicesByCustomer :many
SELECT * FROM invoices
WHERE customer_id = $1
ORDER BY created_at DESC;

-- name: ListInvoicesByOperator :many
SELECT * FROM invoices
WHERE operator_id = $1
ORDER BY created_at DESC;

-- name: UpdateInvoiceStatus :one
UPDATE invoices
SET status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateCustomerPeriodStart :exec
UPDATE customers
SET period_start = $2,
    updated_at = NOW()
WHERE id = $1;
