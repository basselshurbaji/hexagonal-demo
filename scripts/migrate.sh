#!/usr/bin/env bash
# Run goose against db/migrations: migrate.sh <up|down|status|...>
set -euo pipefail
cd "$(dirname "$0")/.."
. scripts/env.sh

exec $GOOSE -dir db/migrations mysql "$DB_DSN" "$@"
