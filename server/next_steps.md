# Noitrex — Remaining Tasks

## Metering Engine
- [x] Write UpsertUsageAggregate sqlc query
- [x] Run sqlc generate
- [x] Implement processEvent — decode message, build Redis key, INCRBY
- [x] Implement flush — read Redis keys, upsert to usage_aggregates, reset counter
- [x] Wire meteringEngine.Start(ctx) in main.go
- [x] Test: POST /ingest → check Redis counter incremented
- [x] Test: wait 30s → check usage_aggregates row created

## Customers
- [x] Write customers migration
- [x] Write sqlc queries — CreateCustomer, GetCustomer, ListCustomers
- [x] Run sqlc generate
- [x] Customer service — CreateCustomer, GetCustomer, ListCustomers
- [x] Customer handler — POST /v1/customers, GET /v1/customers, GET /v1/customers/:id
- [ ] GET /v1/customers/:id/usage — read from Redis

## Plans
- [ ] Write plans migration
- [ ] Write sqlc queries — CreatePlan, GetPlan, ListPlans
- [ ] Run sqlc generate
- [ ] Plan service — CreatePlan, GetPlan, ListPlans
- [ ] Plan handler — POST /v1/plans, GET /v1/plans, GET /v1/plans/:id

## Billing Processor
- [ ] Write invoices migration
- [ ] Write sqlc queries — CreateInvoice, GetInvoice, ListInvoices
- [ ] Run sqlc generate
- [ ] Tier computation function — flat_rate, per_unit, tiered
- [ ] Unit test tier computation for every edge case
- [ ] Billing worker — finds customers whose period ends today
- [ ] Billing worker — reads usage_aggregates, applies plan, creates invoice
- [ ] Billing worker — publishes invoice.created to NexusMQ
- [ ] Wire billing worker in main.go with a daily ticker

## Webhooks
- [ ] Write webhook_endpoints migration
- [ ] Write webhook_deliveries migration
- [ ] Webhook dispatcher — subscribes to invoice.created
- [ ] Webhook dispatcher — looks up operator's endpoints
- [ ] Webhook dispatcher — calls NexusRelay to deliver
- [ ] Webhook handler — POST /v1/webhooks/endpoints

## REST API Polish
- [ ] GET /v1/customers/:id/invoices
- [ ] Consistent error responses across all handlers
- [ ] Request logging middleware
- [ ] Rate limiting on /v1/ingest

## Dashboard
- [ ] Next.js project setup
- [ ] Overview page — revenue, active customers, live counter via SSE
- [ ] Customer detail page — usage graph, invoice history
- [ ] Plan builder — tier editor with live preview

## Done
- [x] Operators — register, login, JWT auth
- [x] POST /v1/ingest — validate, idempotency, create event, publish to NexusMQ
- [x] Metering engine skeleton — Start, subscribe, ticker
- [x] Redis setup — Upstash
- [x] NexusMQ broker setup with topics
- [x] usage_aggregates upsert query