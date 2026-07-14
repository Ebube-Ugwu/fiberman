# fiberman-go-backend

`fiberman-go-backend` is the Go replacement for the current Java backend. It is designed to preserve the existing frontend contract while switching the runtime over to `fiber-go-sdk`.

## Current Scope

Implemented:

- `/api/fiber/node-info`
- `/api/fiber/channels`
- `/api/fiber/channels/details`
- `/api/fiber/peers`
- `/api/fiber/invoices`
- `/api/fiber/payments`
- `/api/fiber/payments/status`
- `/api/fiber/history`
- `/api/settings`
- session-based in-memory history
- runtime settings updates
- generated `curl` snippets
- generated Java and Go SDK snippets
- direct invoice QR generation
- automatic QR payloads when invoice creation returns an invoice string

## Run Locally

```bash
cd fiberman-go-backend
go run ./cmd/server
```

Default environment variables:

- `SERVER_PORT=9020`
- `FIBER_NODE_URL=http://127.0.0.1:8227`
- `FIBER_NODE_AUTH_TOKEN=`
- `FIBER_NODE_TIMEOUT_SECONDS=30`
- `FIBER_PLAYGROUND_BASE_URL=http://localhost:9020`

Example:

```bash
cd fiberman-go-backend
FIBER_NODE_URL=http://127.0.0.1:8227 go run ./cmd/server
```

For the current Angular smoke pass, run the Go backend on `9020` and let the frontend proxy forward `/api` there. This avoids conflicting with the Java backend on `9010`.

## Frontend Compatibility

The Go backend keeps the current REST routes and JSON shapes so the Angular frontend can continue using the same `/api/fiber/*` and `/api/settings` endpoints during migration.
