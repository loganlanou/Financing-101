# Financing 101 Tech Stack

## Backend
- Go 1.25
- Echo v4.13.3 web framework
- Handler / service / repository layering with dependency injection via constructors
- Structured logging with `github.com/lmittmann/tint`
- HTML rendering via [Templ](https://templ.guide) components compiled in `web/components`
- UUIDs generated with `github.com/google/uuid`
- JWT support with `github.com/go-jose/go-jose/v3/jwt` for session validation hooks

## Data Tier
- SQLite (pure Go driver `modernc.org/sqlite`) for the application runtime
- `mattn/go-sqlite3` only used for local tooling and Goose migrations
- Database schema migrations tracked in `db/migrations` and executed with Goose via `make migrate` targets
- Type-safe data access generated via SQLC; queries live in `queries/` and output into `internal/database`

## Integrations
- Authentication managed by Clerk (`internal/auth/clerk.go`)
- Payments through Stripe (`internal/payments/stripe.go`)
- Shipping orchestration via ShipStation (`internal/shipping/shipstation.go`)
- Optional transactional messaging through SendGrid (`internal/mail/sendgrid.go`)

## Tooling
- `Makefile` drives builds, code generation (Templ, SQLC), migrations, and linting
- `go test ./...` used for unit and service tests
- Playwright (configured through `package.json`) reserved for future end-to-end validation
- CSS compiled from `assets/styles` into `web/static` using the npm `sass` script

## Deployment
- Binary compiled with `go build -o logans3d ./cmd/web`
- Systemd-ready service manifests in `deploy/` (future work) reference the compiled binary
- SSR-first architecture: handlers prepare view models, services orchestrate business logic, repositories contain SQLC usages
- Configuration handled via environment variables loaded through `internal/config`

