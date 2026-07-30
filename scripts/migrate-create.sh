#!/usr/bin/env bash
# Create a new migration: migrate-create.sh <name>
set -euo pipefail
cd "$(dirname "$0")/.."
. scripts/env.sh

if [ -z "${1:-}" ]; then
    echo "usage: $0 <migration_name>" >&2
    exit 1
fi

exec $GOOSE -dir db/migrations create "$1" sql
