# Go Finance

Go Finance is a monorepo centered around a Go API for user management, financial record tracking, dashboard summaries, role-based access control, and production-style validation/error handling.

This README is the primary onboarding and API reference for the project.

## Additional Docs

- Backend-specific contributor notes: [app/backend/README.md](app/backend/README.md)
- Repo agent instructions: [AGENTS.md](AGENTS.md)
- Backend agent instructions: [app/backend/AGENTS.md](app/backend/AGENTS.md)

## API Quick Reference

Base URL:

```text
http://localhost:8080
```

For local authenticated requests, these headers are useful:

```bash
-H "X-Dev-Auth-User-Id: local-dev-user" \
-H "X-Dev-User-Role: admin"
```

Important:

- `X-Dev-User-Role` helps local auth identify the caller
- most protected routes still require the user to exist in the `users` table
- `POST /api/v1/users/bootstrap` is typically the first authenticated API call in a fresh local database

### System Endpoints

#### `GET /status`

```bash
curl http://localhost:8080/status
```

#### `GET /docs`

```bash
curl http://localhost:8080/docs
```

#### `GET /static/openapi.json`

```bash
curl http://localhost:8080/static/openapi.json
```

### User Endpoints

#### `POST /api/v1/users/bootstrap`

Creates the initial admin user.

```bash
curl -X POST http://localhost:8080/api/v1/users/bootstrap \
  -H "Content-Type: application/json" \
  -H "X-Dev-Auth-User-Id: local-dev-user" \
  -H "X-Dev-User-Role: admin" \
  -d '{
    "email": "admin@example.com",
    "name": "Admin User"
  }'
```

#### `GET /api/v1/users/me`

```bash
curl http://localhost:8080/api/v1/users/me \
  -H "X-Dev-Auth-User-Id: local-dev-user" \
  -H "X-Dev-User-Role: admin"
```

#### `GET /api/v1/users`

Admin only.

```bash
curl http://localhost:8080/api/v1/users \
  -H "X-Dev-Auth-User-Id: local-dev-user" \
  -H "X-Dev-User-Role: admin"
```

#### `GET /api/v1/users/:id`

Admin only.

```bash
curl http://localhost:8080/api/v1/users/11111111-1111-1111-1111-111111111111 \
  -H "X-Dev-Auth-User-Id: local-dev-user" \
  -H "X-Dev-User-Role: admin"
```

#### `POST /api/v1/users`

Admin only.

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "X-Dev-Auth-User-Id: local-dev-user" \
  -H "X-Dev-User-Role: admin" \
  -d '{
    "authUserId": "user_analyst_1",
    "email": "analyst@example.com",
    "name": "Analyst User",
    "role": "analyst",
    "status": "active"
  }'
```

#### `PATCH /api/v1/users/:id`

Admin only.

```bash
curl -X PATCH http://localhost:8080/api/v1/users/11111111-1111-1111-1111-111111111111 \
  -H "Content-Type: application/json" \
  -H "X-Dev-Auth-User-Id: local-dev-user" \
  -H "X-Dev-User-Role: admin" \
  -d '{
    "role": "viewer",
    "status": "active"
  }'
```

### Financial Record Endpoints

#### `GET /api/v1/records`

Viewer and above.

```bash
curl "http://localhost:8080/api/v1/records?type=expense&category=Food&dateFrom=2026-04-01&dateTo=2026-04-30" \
  -H "X-Dev-Auth-User-Id: local-dev-user" \
  -H "X-Dev-User-Role: viewer"
```

#### `GET /api/v1/records/:id`

Viewer and above.

```bash
curl http://localhost:8080/api/v1/records/11111111-1111-1111-1111-111111111111 \
  -H "X-Dev-Auth-User-Id: local-dev-user" \
  -H "X-Dev-User-Role: viewer"
```

#### `POST /api/v1/records`

Analyst and admin.

```bash
curl -X POST http://localhost:8080/api/v1/records \
  -H "Content-Type: application/json" \
  -H "X-Dev-Auth-User-Id: local-dev-user" \
  -H "X-Dev-User-Role: analyst" \
  -d '{
    "amount": "2500.00",
    "type": "income",
    "category": "Salary",
    "date": "2026-04-01",
    "notes": "Monthly salary"
  }'
```

#### `PATCH /api/v1/records/:id`

Analyst and admin.

```bash
curl -X PATCH http://localhost:8080/api/v1/records/11111111-1111-1111-1111-111111111111 \
  -H "Content-Type: application/json" \
  -H "X-Dev-Auth-User-Id: local-dev-user" \
  -H "X-Dev-User-Role: analyst" \
  -d '{
    "category": "Groceries",
    "notes": "Updated note"
  }'
```

#### `DELETE /api/v1/records/:id`

Admin only.

```bash
curl -X DELETE http://localhost:8080/api/v1/records/11111111-1111-1111-1111-111111111111 \
  -H "X-Dev-Auth-User-Id: local-dev-user" \
  -H "X-Dev-User-Role: admin"
```

### Dashboard Endpoint

#### `GET /api/v1/dashboard/summary`

Viewer and above.

```bash
curl "http://localhost:8080/api/v1/dashboard/summary?dateFrom=2026-01-01&dateTo=2026-03-31&trendInterval=monthly&trendPeriods=3&recentLimit=5" \
  -H "X-Dev-Auth-User-Id: local-dev-user" \
  -H "X-Dev-User-Role: viewer"
```

### Role Summary

- `viewer`: `GET /api/v1/records`, `GET /api/v1/records/:id`, `GET /api/v1/dashboard/summary`
- `analyst`: viewer permissions plus `POST /api/v1/records`, `PATCH /api/v1/records/:id`
- `admin`: analyst permissions plus `DELETE /api/v1/records/:id`, `GET /api/v1/users`, `GET /api/v1/users/:id`, `POST /api/v1/users`, `PATCH /api/v1/users/:id`

## Stack

- Backend: Go, Echo, pgx, PostgreSQL
- Auth: Clerk-compatible auth middleware with local development auth shortcuts
- Background infrastructure: Redis, Asynq, Resend, New Relic
- Shared packages: Bun workspace packages for OpenAPI, Zod schemas, and emails

## Repository Layout

```text
.
├── app/backend              # Go API server
├── packages/openapi         # Shared OpenAPI generation
├── packages/zod             # Shared Zod schemas
├── packages/emails          # Shared email templates
├── package.json             # Root workspace scripts
└── turbo.json               # Turbo pipeline config
```

## Quick Start

### 1. Clone the repo

```bash
git clone <your-repo-url> go-finance
cd go-finance
```

### 2. Install workspace tooling

This repo uses Bun for workspace scripts.

```bash
bun install
```

### 3. Prepare backend environment

The backend reads configuration from `BOILERPLATE_*` environment variables and also autoloads `app/backend/.env` if present.

Create `app/backend/.env` with values like:

```bash
BOILERPLATE_PRIMARY_ENV=local

BOILERPLATE_SERVER_PORT=8080
BOILERPLATE_SERVER_READ_TIMEOUT=30
BOILERPLATE_SERVER_WRITE_TIMEOUT=30
BOILERPLATE_SERVER_IDLE_TIMEOUT=60
BOILERPLATE_SERVER_CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080

BOILERPLATE_DATABASE_HOST=localhost
BOILERPLATE_DATABASE_PORT=5432
BOILERPLATE_DATABASE_USER=postgres
BOILERPLATE_DATABASE_PASSWORD=postgres
BOILERPLATE_DATABASE_NAME=go_finance
BOILERPLATE_DATABASE_SSL_MODE=disable
BOILERPLATE_DATABASE_MAX_OPEN_CONNS=25
BOILERPLATE_DATABASE_MAX_IDLE_CONNS=5
BOILERPLATE_DATABASE_CONN_MAX_LIFETIME=300
BOILERPLATE_DATABASE_CONN_MAX_IDLE_TIME=60

BOILERPLATE_AUTH_SECRET_KEY=sk_test_placeholder
BOILERPLATE_REDIS_ADDRESS=localhost:6379
BOILERPLATE_INTEGRATION_RESEND_API_KEY=re_placeholder
```

Notes:

- In `local` mode, the auth middleware can use `X-Dev-Auth-User-Id` and `X-Dev-User-Role`.
- Redis, Clerk, Resend, and New Relic may still be configured even if you are mainly working on the API layer.

### 4. Start PostgreSQL and create the database

You need a running PostgreSQL instance and a database matching `BOILERPLATE_DATABASE_NAME`.

### 5. Apply migrations

From `app/backend`:

```bash
cd app/backend
BOILERPLATE_DB_DSN="postgres://postgres:postgres@localhost:5432/go_finance?sslmode=disable" task migrations:up
```

### 6. Start the backend

```bash
cd app/backend
task run
```

The server will start on `http://localhost:8080` if you kept the sample port.

## Useful Commands

From the repo root:

```bash
bun install
bun run dev
bun run build
bun run lint
bun run typecheck
```

From `app/backend`:

```bash
task run
task tidy
go test ./...
task migrations:new name=add_some_table
BOILERPLATE_DB_DSN="postgres://..." task migrations:up
```

To rebuild the shared Swagger/OpenAPI spec after backend contract changes:

```bash
cd packages/zod && bun run build
cd ../openapi && bun run build && bun run gen
```

## Authentication and Local Development

Protected routes use auth middleware plus an application-level user lookup.

In local mode, you can call authenticated endpoints with headers like:

```text
X-Dev-Auth-User-Id: local-dev-user
X-Dev-User-Role: admin
```

Important:

- `RequireAuth` identifies the caller.
- `RequireActiveUser` loads the provisioned app user from the `users` table.
- For most protected endpoints, the app user must already exist in the database.
- `POST /api/v1/users/bootstrap` is used to create the initial admin user.

## Role Access Model

The backend now uses explicit permission-based access control.

### Viewer

- Can read financial records
- Can access dashboard summaries
- Cannot create or update records
- Cannot delete records
- Cannot manage users

### Analyst

- Can read financial records
- Can access dashboard summaries
- Can create records
- Can update records
- Cannot delete records
- Cannot manage users

### Admin

- Full user management
- Full record management, including delete
- Can access all read and dashboard capabilities

## API Base URL

When running locally:

```text
http://localhost:8080
```

Versioned API routes live under:

```text
/api/v1
```

## API Documentation UI

- `GET /docs`: serves the API docs UI
- `GET /static/openapi.json`: serves the OpenAPI JSON
- `GET /status`: health check

The OpenAPI document is generated from the shared contract layer:

- `packages/zod`: reusable request and response schemas
- `packages/openapi`: ts-rest contracts and OpenAPI generation

Running `bun run gen` inside `packages/openapi` writes the spec to:

- `packages/openapi/openapi.json`
- `app/backend/static/openapi.json`

## Endpoint Reference

### System Endpoints

#### `GET /status`

Health check endpoint.

Response:

- `200 OK` when the app is healthy
- `503 Service Unavailable` when a dependency check fails

Example response:

```json
{
  "status": "healthy",
  "timestamp": "2026-04-02T12:00:00Z",
  "environment": "local",
  "checks": {
    "database": {
      "status": "healthy",
      "response_time": "2ms"
    },
    "redis": {
      "status": "healthy",
      "response_time": "1ms"
    }
  }
}
```

#### `GET /docs`

Serves the API documentation HTML page.

#### `GET /static/openapi.json`

Serves the OpenAPI JSON document used by the docs UI.

### User Endpoints

#### `GET /api/v1/users/me`

Returns the currently authenticated, provisioned app user.

Auth:

- Requires authentication
- Requires active app user

Response:

- `200 OK`
- `401 Unauthorized`
- `403 Forbidden`

#### `POST /api/v1/users/bootstrap`

Creates the initial admin user for the application. This is intended for first-time setup.

Auth:

- Requires authentication
- Does not require an already provisioned app user

Request body:

```json
{
  "email": "admin@example.com",
  "name": "Admin User"
}
```

Response:

- `201 Created`
- `400 Bad Request` for invalid input
- `403 Forbidden` if an initial admin already exists

#### `GET /api/v1/users`

Lists all users.

Auth:

- Requires authentication
- Requires active app user
- Requires admin permission

Response:

- `200 OK`
- `401 Unauthorized`
- `403 Forbidden`

#### `GET /api/v1/users/:id`

Fetches a single user by UUID.

Auth:

- Admin only

Response:

- `200 OK`
- `404 Not Found`

#### `POST /api/v1/users`

Creates a user.

Auth:

- Admin only

Request body:

```json
{
  "authUserId": "user_123",
  "email": "viewer@example.com",
  "name": "Viewer User",
  "role": "viewer",
  "status": "active"
}
```

Response:

- `201 Created`
- `400 Bad Request`
- `403 Forbidden`

#### `PATCH /api/v1/users/:id`

Updates an existing user.

Auth:

- Admin only

Rules:

- At least one field must be provided
- `role` must be one of `viewer`, `analyst`, `admin`
- `status` must be one of `active`, `inactive`

Example body:

```json
{
  "role": "analyst",
  "status": "active"
}
```

Response:

- `200 OK`
- `400 Bad Request`
- `404 Not Found`

### Financial Record Endpoints

#### `GET /api/v1/records`

Lists financial records.

Auth:

- Viewer and above

Query params:

- `type`: `income` or `expense`
- `category`: exact category filter
- `dateFrom`: `YYYY-MM-DD`
- `dateTo`: `YYYY-MM-DD`

Response:

- `200 OK`
- `400 Bad Request`

#### `GET /api/v1/records/:id`

Returns a single financial record by UUID.

Auth:

- Viewer and above

Response:

- `200 OK`
- `404 Not Found`

#### `POST /api/v1/records`

Creates a financial record.

Auth:

- Analyst and admin

Request body:

```json
{
  "amount": "2500.00",
  "type": "income",
  "category": "Salary",
  "date": "2026-04-01",
  "notes": "April salary"
}
```

Rules:

- `amount` must be a valid non-negative decimal
- `type` must be `income` or `expense`
- `category` is required
- `date` must be `YYYY-MM-DD`

Response:

- `201 Created`
- `400 Bad Request`

#### `PATCH /api/v1/records/:id`

Updates a financial record.

Auth:

- Analyst and admin

Rules:

- At least one field must be provided
- Field validations are the same as create

Response:

- `200 OK`
- `400 Bad Request`
- `404 Not Found`

#### `DELETE /api/v1/records/:id`

Deletes a financial record.

Auth:

- Admin only

Response:

- `204 No Content`
- `404 Not Found`

### Dashboard Summary Endpoint

#### `GET /api/v1/dashboard/summary`

Returns aggregated dashboard data for the authenticated user.

Auth:

- Viewer and above

Query params:

- `dateFrom`: optional `YYYY-MM-DD`
- `dateTo`: optional `YYYY-MM-DD`
- `trendInterval`: optional `weekly` or `monthly`
- `trendPeriods`: optional integer `1-24`
- `recentLimit`: optional integer `1-20`

Response includes:

- `totalIncome`
- `totalExpenses`
- `netBalance`
- `categoryTotals`
- `recentActivity`
- `trends`

Example response:

```json
{
  "dateFrom": "2026-01-01T00:00:00Z",
  "dateTo": "2026-03-31T00:00:00Z",
  "trendInterval": "monthly",
  "totalIncome": "6000.00",
  "totalExpenses": "2400.00",
  "netBalance": "3600.00",
  "categoryTotals": [
    {
      "type": "expense",
      "category": "Rent",
      "total": "1500.00"
    }
  ],
  "recentActivity": [],
  "trends": [
    {
      "interval": "monthly",
      "periodStart": "2026-01-01T00:00:00Z",
      "periodEnd": "2026-01-31T00:00:00Z",
      "income": "2000.00",
      "expenses": "800.00",
      "netBalance": "1200.00"
    }
  ]
}
```

Response:

- `200 OK`
- `400 Bad Request`

## Validation and Error Handling

The backend demonstrates production-style request handling:

- Struct validation for required and constrained fields
- Strict JSON decoding for write endpoints
- Rejection of malformed JSON
- Rejection of unknown JSON fields
- Rejection of invalid types
- Field-level validation errors in a structured response
- Appropriate HTTP status codes for validation, auth, forbidden, not found, and server failures

Typical error response shape:

```json
{
  "code": "BAD_REQUEST",
  "message": "Validation failed",
  "status": 400,
  "override": true,
  "errors": [
    {
      "field": "dateFrom",
      "error": "must be before or equal to dateTo"
    }
  ],
  "action": null
}
```

Common statuses:

- `400 Bad Request`: invalid or incomplete input
- `401 Unauthorized`: missing or invalid auth
- `403 Forbidden`: authenticated but not allowed
- `404 Not Found`: resource or route not found
- `429 Too Many Requests`: rate limited
- `500 Internal Server Error`: unexpected failure
- `503 Service Unavailable`: unhealthy service dependencies

## Testing

Run backend tests with:

```bash
cd app/backend
go test ./...
```

Recent backend work includes tests for:

- Service-level financial record logic
- Auth and authorization behavior
- Validation and malformed payload handling

## Notes for Contributors

- Keep handlers thin and push business logic into services
- Keep repository code focused on data access
- Preserve the explicit permission-based authorization model
- Add tests alongside any new middleware, handler, or service behavior
- Keep the OpenAPI output in sync with runtime endpoints
