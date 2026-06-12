# Noitrex — Remaining Tasks

## Metering Engine
- [xL] Write UpsertUsageAggregate sqlc query
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
- [ ] GET /v1/customers/:id/usage — read from Redis - (will be implemented once all core api's and endpoints have been implemented.)

## Plans
- [x] Write plans migration
- [x] Write sqlc queries — CreatePlan, GetPlan, ListPlans
- [x] Run sqlc generate
- [x] Plan service — CreatePlan, GetPlan, ListPlans
- [x] Plan handler — POST /v1/plans, GET /v1/plans, GET /v1/plans/:id

## Billing Processor
- [x] Write invoices migration
- [x] Write sqlc queries — CreateInvoice, GetInvoice, ListInvoices
- [x] Run sqlc generate
- [x] Tier computation function — flat_rate, per_unit, tiered
- [ ] Unit test tier computation for every edge case
- [x] Billing worker — finds customers whose period ends today
- [x] Billing worker — reads usage_aggregates, applies plan, creates invoice
- [x] Billing worker — publishes invoice.created to NexusMQ
- [x] Wire billing worker in main.go with a daily ticker

## Webhooks
- [x] Write webhook_endpoints migration
- [x] Write webhook_deliveries migration
- [x] Webhook dispatcher — subscribes to invoice.created
- [x] Webhook dispatcher — looks up operator's endpoints
- [x] Webhook dispatcher — calls NexusRelay to deliver
- [x] Webhook handler — POST /v1/webhooks/endpoints

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