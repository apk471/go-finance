# Backend Guidelines

## Scope

These instructions apply to `app/backend`.

## Architecture

The backend follows a layered Go structure:

- `internal/handler`: HTTP request and response handling
- `internal/service`: business logic and validation beyond transport concerns
- `internal/repository`: database access
- `internal/middleware`: auth, permissions, request context, global error handling
- `internal/model`: shared domain models and response DTOs

Keep new code consistent with this structure.

## Existing Feature Areas

The current backend includes:

- Health and docs endpoints
- User bootstrap and admin user management
- Financial record listing, creation, update, and deletion
- Dashboard summary aggregation
- Permission-based access control
- Structured validation and structured HTTP error responses

Prefer extending these existing patterns instead of creating alternate flows.

## Route and Auth Patterns

- Public/system routes are registered from `internal/router/system.go`
- Versioned API routes live under `/api/v1`
- Protected routes should typically use:
  - `RequireAuth`
  - `RequireActiveUser`
  - `RequirePermission(...)`

When adding a new protected route, decide the exact permission instead of relying on incidental route grouping.

## Validation Patterns

- Request structs should implement `Validate() error`
- Use validator tags for straightforward field rules
- Use `validation.CustomValidationErrors` for cross-field or request-level checks
- For write endpoints, keep the current strict JSON behavior intact

## Error Handling Patterns

- Return `errs.HTTPError` for expected user-facing failures
- Let the global error handler format responses consistently
- Use meaningful status codes:
  - `400` for invalid input
  - `401` for missing or invalid auth
  - `403` for blocked actions
  - `404` for missing resources
  - `500` for unexpected failures

## Testing

At minimum for backend changes:

- run `gofmt` on touched Go files
- run `cd app/backend && go test ./...`

Add tests for:

- new service rules
- authorization decisions
- validation failures
- changed handler behavior

## Operations and Docs

- Use `task run` to start the backend locally
- Use `task migrations:new` and `task migrations:up` for schema changes
- Keep `README.md`, `app/backend/README.md`, and runtime behavior in sync
- Keep `app/backend/static/openapi.json` aligned with the exposed API where applicable
