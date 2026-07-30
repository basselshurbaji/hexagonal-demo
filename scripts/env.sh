# Sourced by the other scripts: loads .env and derives DB_DSN.
# Defaults mirror config.go so scripts work even without a .env file.
set -a
[ -f .env ] && . ./.env
set +a

: "${DB_HOST:=mysql.hexagonal-demo.orb.local}"
: "${DB_PORT:=3306}"
: "${DB_DATABASE:=hexagonal}"
: "${DB_USERNAME:=app}"
: "${DB_PASSWORD:=app}"

DB_DSN="${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}?parseTime=true"
GOOSE="go run github.com/pressly/goose/v3/cmd/goose@v3.26.0"
