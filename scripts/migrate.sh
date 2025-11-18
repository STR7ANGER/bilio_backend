#!/usr/bin/env bash
set -euo pipefail

if ! command -v psql >/dev/null 2>&1; then
  echo "[error] psql is required to run migrations"
  exit 1
fi

DATABASE_URL=${DATABASE_URL:-"postgresql://user:password@localhost:5432/bilio"}
MIGRATIONS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/migrations"

if [ ! -d "$MIGRATIONS_DIR" ]; then
  echo "[error] migrations directory not found: $MIGRATIONS_DIR"
  exit 1
fi

shopt -s nullglob
readarray -t MIGRATIONS < <(printf '%s\n' "$MIGRATIONS_DIR"/*.sql | sort)
shopt -u nullglob

if [ ${#MIGRATIONS[@]} -eq 0 ]; then
  echo "[warn] no migrations to apply"
  exit 0
fi

for migration in "${MIGRATIONS[@]}"; do
  echo "[info] applying $(basename "$migration")"
  psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f "$migration"
done

echo "[info] migrations applied"
