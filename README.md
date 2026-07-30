# hexagonal-demo

A demo / proof of concept of [Hexagonal Architecture (Ports & Adapters)](https://alistair.cockburn.us/hexagonal-architecture/)
in pure Go with minimal dependencies — a modular monolith of three modules
(`user`, `song`, `playlist`) over MySQL.

Runtime dependencies: the MySQL driver and `godotenv`. Tooling
([goose](https://github.com/pressly/goose) for migrations,
[sqlc](https://sqlc.dev) for type-safe queries) runs via `go run` — nothing to
install.

## Architecture

Each module is its own hexagon:

```
modules/<name>/
├── <name>.go     Module + Initialize()   — assembly from driven ports
├── facade.go     public API: models, methods, exported errors
├── adapters/     driving (http) and driven (sql, other modules)
└── internal/     the hexagon — inaccessible outside the module
    ├── entity/   domain models
    ├── ports/    interfaces the core depends on
    └── service/  use-cases
```

Rules the code follows:

- **Dependencies point inward.** `entity`, `ports`, `service` import nothing
  from adapters, sqlc, or `net/http`. Adapters depend on the core, never the
  reverse.
- **The facade is the module boundary.** Other modules and driving adapters
  consume public models and errors from the module root package; internal
  entities never leave the hexagon.
- **Modules communicate through consumer-owned ports.** e.g. the playlist
  core defines `ports.SongCatalog` in its own vocabulary; an adapter
  translates the song module's facade into it. Swapping that adapter for an
  HTTP client would extract the module into a service without touching the core.
- **No module touches another module's tables.** Data crosses module
  boundaries through ports only.
- **`main` is the composition root.** `register.go` wires everything in two
  phases: create all module references, then initialize them — so dependency
  cycles between modules cannot block construction.

## Quick start

```sh
cp .env.example .env
make up        # start MySQL (docker compose, waits for healthy)
make migrate   # apply goose migrations
make seed      # seed demo data (20 artists, 100 songs, 1 user)
go run .
```

## Endpoints

| Method | Path                        | Description                          |
|--------|-----------------------------|--------------------------------------|
| GET    | `/health`                   | liveness + DB ping                   |
| GET    | `/users/{id}`               | get user                             |
| GET    | `/users/{id}/playlists`     | playlists of a user                  |
| GET    | `/songs?ids=1,2,3`          | songs by ids                         |
| POST   | `/playlists`                | create playlist `{name, user_id}`    |
| PUT    | `/playlists/{id}/songs`     | add songs `{song_ids}` (validated)   |
| GET    | `/playlists/{id}`           | playlist with its songs              |

## Database workflow

- `make migrate-create name=...` — new goose migration in `db/migrations`
- `make migrate` / `migrate-down` / `migrate-status`
- `make seed` / `seed-reset` — goose `-no-versioning` over `db/seeds`
- `make sqlc` — regenerate `db/gen` from migrations (schema) + `db/queries`;
  the generated code is committed, so rerun this after changing queries or schema
