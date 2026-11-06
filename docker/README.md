# Go API Docker stack

This Compose stack runs the Go API, the Nuxt admin app, PostgreSQL, and Redis.

## Start

```powershell
Copy-Item docker\api.toml.example docker\api.toml
docker compose up --build -d
```

Open <http://localhost:23324> for the API and <http://localhost:23325> for the admin app. The web installer has been removed: run `pnpm dev:db` from the repository root to create the tables and the first administrator. That command reads the root `config.toml`, so point it at `127.0.0.1:5433` (database and user `forestwhisper`, password `forestwhisper-dev`) before running it. Docker-side initialization is not wired up yet.

`docker/api.toml` is the API's runtime config and is gitignored. It already contains the Compose connection settings, so no installer step is needed. PostgreSQL uses the same local-development-only password in Compose; change it in both files before starting a fresh stack if needed. This committed default is for local development, not production.

The root `config.toml` is not mounted or changed. PostgreSQL and Redis data persist in named volumes.

## What the Docker files do

- `go.Dockerfile` builds the server from `apps/api/cmd/servers`, then copies only the binary into a small non-root runtime image.
- `nuxt.Dockerfile` builds the Nuxt admin app (`apps/admin`) with pnpm and runs its Nitro node server on port `23325`; there is no web server inside the container.
- Compose connects containers by service name on its default network (`postgres:5432`, `redis:6379`, `forwhis-api:23324`); the host-only ports are `23324`, `23325`, `5433`, and `6380`.
- API database and Redis settings come from `docker/api.toml`; no `.env` file is used.
- Database health checks gate API startup. Redis append-only persistence and named volumes keep service data across container recreation.

## Admin app

`docker/nuxt.Dockerfile` builds the app inside the image: it installs the pnpm workspace for `apps/admin`, runs `nuxt build` on Node 26, then copies only Nitro's `.output` into a clean Node runtime image that starts `node .output/server/index.mjs` on port `23325`. Nothing is built on the host, so the artifacts always match the runtime platform. There is no web server inside the container; Nitro serves client-side routes itself, so a reverse proxy in front of it only has to forward traffic.

`ADMIN_API_BASE_URL` sets `NUXT_PUBLIC_API_BASE_URL`, both at build time and as a container env var, and defaults to `/` (same origin as the admin). It is the address the browser calls, so it has to be publicly reachable:

```powershell
$env:ADMIN_API_BASE_URL = "https://api.example.com"; docker compose up -d --build admin
```

```powershell
docker compose logs -f forwhis-api
docker compose down
```

`docker compose down` keeps the named database volumes and `docker/api.toml`. `docker compose down -v` deletes PostgreSQL and Redis data.
