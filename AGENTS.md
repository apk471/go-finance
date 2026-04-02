# Repository Guidelines

## Project Structure

This monorepo contains a Go backend plus shared TypeScript packages.

- `app/backend/`: main API server
- `packages/openapi`: shared OpenAPI generation
- `packages/zod`: shared validation schemas
- `packages/emails`: shared email templates
- `README.md`: primary onboarding and API documentation

## Development Commands

Run from the repo root unless a command says otherwise.

- `bun install`: install workspace dependencies
- `bun run dev`: start workspace dev pipelines
- `bun run build`: build all configured workspaces
- `bun run lint`: run lint tasks across workspaces
- `bun run typecheck`: run TypeScript type checks
- `cd app/backend && task run`: start the Go API
- `cd app/backend && go test ./...`: run backend tests
- `cd app/backend && task tidy`: format and tidy backend dependencies
- `cd app/backend && task migrations:new name=add_feature`: create a migration
- `cd app/backend && BOILERPLATE_DB_DSN=... task migrations:up`: apply migrations

## Backend Conventions

- Keep the Go backend layered as `handler -> service -> repository`
- Keep handlers thin and push business rules into services
- Keep repository methods focused on data access and SQL concerns
- Use `gofmt` on all edited Go files
- Add tests next to the changed package as `*_test.go`

## Current Domain

The backend currently supports:

- User bootstrap and user management
- Financial record CRUD
- Dashboard summary aggregation APIs
- Permission-based access control for `viewer`, `analyst`, and `admin`
- Structured validation and error handling

When adding new features, preserve those patterns instead of introducing parallel ones.

## Access Control Expectations

The repo now uses explicit permissions instead of relying only on role ordering.

- `viewer`: read records and summaries
- `analyst`: read plus create and update records
- `admin`: full record and user management

If you add protected endpoints, wire them through the authorization middleware and follow the existing permission model.

## Validation and Errors

- Prefer request structs with `Validate()` methods
- Return structured `errs.HTTPError` responses for user-facing failures
- Use field-level validation errors for bad input
- Reject invalid operations early in the service layer when possible

## Documentation Expectations

- Keep `README.md` aligned with actual endpoints and startup steps
- Keep `app/backend/README.md` aligned with backend-specific workflows
- Update both `AGENTS.md` files when repo conventions materially change

## Security and Config

- Never commit secrets
- Use `BOILERPLATE_*` environment variables for backend config
- Set `BOILERPLATE_DB_DSN` explicitly for migration commands
- Keep generated docs in sync with `app/backend/static/openapi.json`
