-- name: UpsertUsageAggregate :exec
INSERT INTO usage_aggregates
(operator_id,customer_id,event_name,period_start,period_end,total_units)
VALUES ($1,$2,$3,$4,$5,$6)
ON CONFLICT (customer_id,event_name,period_start) DO UPDATE SET total_units = usage_aggregates.total_units + EXCLUDED.total_units, updated_at = NOW();