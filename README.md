# Bilio Backend

A starter Go backend with a lightweight PostgreSQL data layer, organized for clean architecture and ready for containerization or deployment.

## Getting Started

1. **Install dependencies**
   ```bash
   go mod tidy
   ```

2. **Set up environment**
   ```bash
   cp .env.example .env
   ```
   Configure `DATABASE_URL`, `EMAIL_USER`, and `EMAIL_PASSWORD` for your SMTP provider (e.g. Gmail, SendGrid).

3. **Run database migrations (optional)**
   The server automatically ensures required tables exist on startup. If you prefer to run the SQL ahead of time, execute:
   ```bash
   ./scripts/migrate.sh
   ```

4. **Start the development server**
   ```bash
   ./scripts/dev.sh
   ```

## Project Structure

- `cmd/server` — application entrypoint
- `internal/app` — domain-specific logic (handlers, services, repositories, models)
- `internal/config` — configuration loading
- `internal/database` — PostgreSQL client bootstrap and schema helpers
- `internal/logger` — structured logging
- `internal/server` — HTTP server wiring
- `pkg/mailer` — outbound email integrations
- `pkg/middleware` — reusable middleware
- `scripts` — helper scripts for development and tooling

## Database

The application connects to PostgreSQL using the `pgx` driver. Connection details come from `DATABASE_URL` (defaulting to `postgresql://user:password@localhost:5432/bilio`).

On startup the service will ensure the `users` and `waitlist_entries` tables exist. For production workloads you should manage migrations explicitly (e.g. via Goose, Flyway, or another tool) and disable automatic schema creation.

## Email

The service sends waitlist emails through SMTP. Provide credentials via the following environment variables:

- `EMAIL_USER` — SMTP username/login
- `EMAIL_PASSWORD` — SMTP password or app-specific token
- `APP_EMAIL_FROM` (optional) — overrides the default `waitlist@billstack.com`

If these variables are missing the server will fail to boot.

## Hot Reloading

This project ships with an `.air.toml` configuration for [Air](https://github.com/air-verse/air). Install Air (e.g. `go install github.com/air-verse/air@latest`) and run `./scripts/dev.sh` to start the server with automatic reloads.

## Next Steps

- Implement application-specific models and handlers.
- Add request validation, authentication, and testing.
- Containerize the service or integrate with CI/CD.

