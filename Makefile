.PHONY: up down migrate migrate-down migrate-status migrate-create seed seed-reset sqlc

## up: start services (waits for MySQL healthcheck)
up:
	docker compose up -d --wait

## down: stop services (data volume is kept; use `docker compose down -v` to wipe)
down:
	docker compose down

## migrate: apply all pending migrations
migrate:
	./scripts/migrate.sh up

## migrate-down: roll back the last migration
migrate-down:
	./scripts/migrate.sh down

## migrate-status: show migration status
migrate-status:
	./scripts/migrate.sh status

## migrate-create: create a new migration, e.g. `make migrate-create name=create_users`
migrate-create:
	./scripts/migrate-create.sh $(name)

## seed: load seed data (errors on duplicates if already seeded — run seed-reset first to reapply)
seed:
	./scripts/seed.sh up

## seed-reset: remove seed data (runs the Down of every seed file)
seed-reset:
	./scripts/seed.sh down-to 0

## sqlc: generate Go code from schema + queries (cleans first: sqlc doesn't remove outputs of renamed/deleted query files)
sqlc:
	rm -rf db/gen
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0 generate
