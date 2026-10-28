# Go API Docker stack

This Compose stack runs only the Go API, PostgreSQL, and Redis. The Vue/Web project is not built or started here.

## Start

```powershell
Copy-Item docker\api.toml.example docker\api.toml
docker compose up --build -d
```

Open <http://localhost:2324> and complete the first-run installer. Enter these PostgreSQL values (replace the install page's default `127.0.0.1`; inside the API container, PostgreSQL is `postgres`):

- Host: `postgres`
- Port: `5432`
- Database and user: `forestwhisper`
- Password: `forestwhisper-dev`

`docker/api.toml` is the API's writable runtime config and is gitignored. The installer saves its settings there. PostgreSQL uses the same local-development-only password in Compose; change it in both files before starting a fresh stack if needed. This committed default is for local development, not production.

The root `config.toml` is not mounted or changed. PostgreSQL and Redis data persist in named volumes.

## What the Docker files do

- `go.Dockerfile` builds the server from `apps/api/cmd/servers`, then copies only the binary into a small non-root runtime image.
- Compose connects containers by service name on its default network (`postgres:5432`, `redis:6379`); the host-only ports are `2324`, `5433`, and `6380`.
- API database and Redis settings come from `docker/api.toml`; no `.env` file is used.
- Database health checks gate API startup. Redis append-only persistence and named volumes keep service data across container recreation.

```powershell
docker compose logs -f forwhis-api
docker compose down
```

`docker compose down` keeps the named database volumes and `docker/api.toml`. `docker compose down -v` deletes PostgreSQL and Redis data.
