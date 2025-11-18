#!/usr/bin/env bash
set -euo pipefail

if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

if ! command -v air >/dev/null 2>&1; then
  echo "[warn] air not installed; falling back to go run"
  go run ./cmd/server
else
  air
fi
