# Noitrex

Noitrex is a self-hostable, usage-based billing engine for SaaS APIs built in Go. It is a highly performant B2B infrastructure tool designed to count usage, compute invoices, and fire webhooks — without getting in the way.

## 🤝 The Three-Party Model

Noitrex operates on a strict three-party relationship:
1. **Operator**: A SaaS company (e.g., an AI API startup) that self-hosts or integrates Noitrex to handle billing for their own customers. Operators authenticate via API keys for server-to-server requests and via JWT for the dashboard.
2. **Customer**: The operator's end-user. Customers are billed by the operator based on their usage.
3. **Noitrex**: The engine in the middle. It counts usage, computes invoices based on pricing plans, and fires webhook notifications back to the Operator.

## 🔄 Core Data Flow

Noitrex is designed for speed and reliability, utilizing a hybrid persistence model:

1. **Ingestion**: The Operator sends usage events to `POST /ingest` in real-time. Each event represents a billable action performed by a Customer and includes an idempotency key to guarantee exactly-once counting.
2. **Persistence & Internal Pub/Sub**: Noitrex immediately persists the raw event to PostgreSQL, publishes it internally via **NexusMQ** (an in-process pub/sub broker), and returns an HTTP 200.
3. **Fast Counting**: The metering engine subscribes to NexusMQ ingestion events and increments a Redis counter per customer, per billing period.
4. **Aggregate Flushing**: Redis counters flush their states to PostgreSQL aggregates on a strict 30-second interval using atomic `INSERT ... ON CONFLICT DO UPDATE` queries.
5. **Invoice Generation**: At the end of a billing period, a background worker reads the PostgreSQL aggregates, applies the operator's pricing plan, computes the invoice amount in **integer kobo** (never floats), and writes the immutable invoice record.
6. **Webhook Dispatch**: Invoice creation triggers a NexusMQ event. The webhook dispatcher picks this up and delivers a signed HMAC-SHA256 webhook to the operator's registered endpoint via **NexusRelay** (our dedicated webhook delivery library).

## 💰 Pricing Models Supported

- **Flat Rate**: Fixed price per unit regardless of volume.
- **Per Unit**: Linear pricing with an optional free tier.
- **Tiered**: Usage walks through price buckets sequentially. Each tier only prices the units that fall within its specific range.

## 🧱 Technology Stack

- **Backend**: Go, Gin
- **Database**: PostgreSQL (pgx, sqlc, goose)
- **Caching & Real-time Counters**: Redis
- **Internal Messaging**: NexusMQ (In-process Pub/Sub)
- **Webhook Delivery**: NexusRelay
- **Dashboard**: Next.js App Router (React)

## 🚫 What Noitrex Does NOT Do

To maintain strict boundaries and extreme focus, Noitrex explicitly avoids the following:

- **No Card Charging**: Noitrex computes what is owed and fires a webhook. The Operator is entirely responsible for payment collection via Stripe, Paystack, or their provider of choice.
- **No External Message Queues**: NexusMQ is strictly in-process for v1. We do not use Kafka or Redis Streams.
- **No Cross-Instance Multi-Tenancy**: One Noitrex deployment serves one operator group. Operator isolation is strictly enforced at the database level via `operator_id` scoping on every query.

## ⚙️ Key Engineering Constraints

- **Money is Integer Kobo**: All monetary values are stored as `BIGINT` representing the smallest currency unit. Floats are strictly prohibited anywhere near financial calculations.
- **Database-Level Idempotency**: Idempotency keys have a `UNIQUE` constraint at the database schema level. Duplicate events are rejected natively at the insert step, bypassing application-layer race conditions.
- **Atomic Increments**: Usage aggregates utilize `INSERT ... ON CONFLICT DO UPDATE SET total_units = total_units + EXCLUDED.total_units` to prevent read-modify-write race conditions.
- **Immutable Invoices**: Once an invoice is issued, it is strictly immutable. Any corrections must be modeled via new "Adjustment" records, never by `UPDATE`ing the invoice.

---
*Noitrex: Stop building billing infra. Start building your product.*
