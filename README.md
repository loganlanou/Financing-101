# Financing 101 (Go + Echo + Templ)

A ground-up rebuild of the Investing101 prototype using the requested production stack:
Go 1.25, Echo v4.13.3, Templ SSR components, SQLC-powered repositories, and Goose-run
SQLite migrations. The app mirrors the original feature set—news + NLP sentiment,
S&P 500 benchmarking, congressional trade tracking, and curated AI recommendations—while
layering handler/service/repository boundaries to prep for the future mobile app.

## Getting Started

### Prerequisites
- Go 1.25+
- Node.js 20+ (Sass + Playwright tooling)
- SQLite tooling (optional, for inspecting `data/app.db`)

### Install deps
```bash
make deps   # optional alias: go mod tidy && npm install
npm install
```

### Generate + migrate
```bash
make generate   # templ + sqlc codegen
make migrate-up # goose migrations -> data/app.db
```

### Seed demo data
Seeded automatically on boot, or run manually:
```bash
go run ./cmd/tools/seed   # placeholder for future CLI
```

### Build & Run
```bash
make build
./logans3d
# or
go run ./cmd/web
```
Visit `http://localhost:8080`.

### GitHub Pages Preview
`docs/index.html` + `docs/app.css` provide a static snapshot of the dashboard for
GitHub Pages. Point Pages at the `docs/` folder (Settings → Pages → Build and deployment → Source: Deploy from branch → Branch: `main`, Folder: `/docs`).  
After the next push, `https://<user>.github.io/Financing-101/` will render the live-looking preview instead
of the README.

### Configuration
- `NEWS_FEEDS`: comma-separated RSS feeds (defaults to Yahoo Finance + WSJ Markets).
- `NEWS_POLL_INTERVAL`: cadence for refreshing feeds (default `30m`).
- `REQUEST_TIMEOUT`: guards handler + ingest calls (default `4s`).
- Standard credentials: `CLERK_SECRET_KEY`, `STRIPE_SECRET_KEY`, `SHIPSTATION_*`, `SENDGRID_API_KEY`, `SIGNING_KEY`, etc.

### CSS workflow
```bash
npm run build:css   # one-off Sass compile
npm run watch:css   # dev loop
```

### Tests
```bash
go test ./...
npm run test        # Playwright test placeholder
```

## Project Structure
```
.
├── cmd/web                 # main entrypoint + Goose wiring
├── db/migrations           # Goose SQL + embedded FS package
├── internal/
│   ├── auth|mail|payments|shipping stubs
│   ├── config/logging
│   ├── data                # seed orchestrator
│   ├── handlers            # Echo handlers -> templ components
│   └── services            # business logic on top of SQLC queries
├── queries/                # SQLC query definitions
├── web/components          # templ SSR components
├── web/static              # compiled CSS output
├── assets/styles           # Sass source
├── Makefile                # build/generate/migrate targets
└── package.json            # Sass + Playwright infra
```

## Feature Parity Highlights
- **News + Sentiment**: multi-source RSS ingestion with govader sentiment scoring and ticker extraction.
- **Stock Lab**: three timeframes of performance plus vs S&P delta (mirrors CLI prototype).
- **Congress Watch**: top disclosures (Nancy Pelosi, etc.) with sentiment heat.
- **AI Desk**: aggregated recommendations w/ conviction scoring.
- **Integrations**: Clerk (JWT verification via go-jose), Stripe, ShipStation, SendGrid stubs wired for future implementations.
- **Tooling**: Goose migrations, SQLC repositories, templ SSR, Sass build via npm, Playwright placeholder.

## Deployment Notes
- `make build` outputs `logans3d` binary (staging/prod systemd target).
- Provide required env vars (`HTTP_ADDR`, `DATABASE_PATH`, `CLERK_SECRET_KEY`, etc.) via systemd drop-ins or `.env`.
- For migrations in CI/CD, run `make migrate-up` (Goose uses embedded FS ensuring binaries carry schema).
