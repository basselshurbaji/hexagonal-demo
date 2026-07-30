#!/usr/bin/env bash
# Run goose against db/seeds (not version-tracked): seed.sh <up|down-to 0|...>
set -euo pipefail
cd "$(dirname "$0")/.."
. scripts/env.sh

exec $GOOSE -dir db/seeds -no-versioning mysql "$DB_DSN" "$@"
