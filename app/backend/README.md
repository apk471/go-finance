# Backend README

This file is the backend-focused companion to the root [README.md](/Users/ayushamin/Developer/repos/go-finance/README.md).

## What Lives Here

`app/backend` contains the Go API server, including:

- Echo router and handlers
- auth and authorization middleware
- service and repository layers
- PostgreSQL migrations
- static API docs assets
- email templates

## Common Commands

Run from `app/backend`:

```bash
task run
task tidy
go test ./...
task migrations:new name=add_some_change
BOILERPLATE_DB_DSN="postgres://..." task migrations:up
```

## Config

The backend loads config from `BOILERPLATE_*` environment variables and also autoloads `.env`.

Required config groups in practice:

- `primary`
- `server`
- `database`
- `auth`
- `redis`
- `integration`

If any required values are missing, startup fails during config validation.

## API Surface

Current backend route groups:

- System:
  - `GET /status`
  - `GET /docs`
  - `GET /static/openapi.json`
- Users:
  - `GET /api/v1/users/me`
  - `POST /api/v1/users/bootstrap`
  - `GET /api/v1/users`
  - `GET /api/v1/users/:id`
  - `POST /api/v1/users`
  - `PATCH /api/v1/users/:id`
- Financial records:
  - `GET /api/v1/records`
  - `GET /api/v1/records/:id`
  - `POST /api/v1/records`
  - `PATCH /api/v1/records/:id`
  - `DELETE /api/v1/records/:id`
- Dashboard:
  - `GET /api/v1/dashboard/summary`

See the root README for request and response details.

## OpenAPI Source Of Truth

The backend docs UI serves `static/openapi.json`, but that file is generated from the shared packages:

- `packages/zod`: request and response schemas
- `packages/openapi`: route contracts and OpenAPI generation

To refresh the backend spec after API changes:

```bash
cd packages/zod && bun run build
cd ../openapi && bun run build && bun run gen
```

## Authorization Model

Permissions are enforced in middleware.

- `viewer`: read records and summaries
- `analyst`: read plus create and update records
- `admin`: delete records and manage users

If you add a new protected route, use `RequirePermission(...)` and extend the permission map if needed.

## Validation and Errors

The backend uses:

- request structs with `Validate()`
- strict JSON decoding for write requests
- structured `errs.HTTPError` responses
- SQL error mapping through the global error handler

This means malformed JSON, unknown fields, invalid UUIDs, missing required fields, and blocked actions should produce consistent API errors instead of raw framework errors.

## Testing Guidance

Keep tests near the code they cover.

- middleware tests in `internal/middleware`
- handler tests in `internal/handler`
- service tests in `internal/service`

Recent work especially relies on tests for:

- financial summary logic
- authorization decisions
- validation failures and malformed payload handling
