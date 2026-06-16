# Deployment guide

## Local full-stack with Docker Compose

Prerequisites: Docker Engine 24+ and Docker Compose v2.

```bash
docker compose up --build
```

| Service | URL |
|---------|-----|
| Web frontend | http://localhost |
| API | http://localhost/api/v1/healthz |
| Readiness probe | http://localhost/api/v1/readyz |
| Metrics | http://localhost/debug/vars |

The first startup runs schema migrations automatically against the
Postgres container. Case data is stored in the `db_data` named volume
and survives container restarts.

## Environment variables

### API server (`faraidd`)

| Variable | Default | Description |
|----------|---------|-------------|
| `FARAID_ENV` | `development` | One of `development`, `test`, `production`. Production requires `FARAID_DATABASE_URL`. |
| `FARAID_HTTP_ADDR` | `:8080` | Host:port the HTTP server listens on. |
| `FARAID_LOG_LEVEL` | `info` | One of `debug`, `info`, `warn`, `error`. |
| `FARAID_LOG_FORMAT` | `json` | One of `json`, `text`. |
| `FARAID_DATABASE_URL` | _(empty)_ | PostgreSQL DSN. When unset the in-memory store is used (cases are lost on restart). |

### Web server (SvelteKit node)

| Variable | Default | Description |
|----------|---------|-------------|
| `ORIGIN` | _(required in production)_ | The public URL the server is accessed from, e.g. `https://faraid.example.com`. Needed by the CSRF protection built into SvelteKit. |
| `PORT` | `3000` | Port the Node.js server listens on. |
| `PUBLIC_LLM_ENABLED` | _(empty)_ | Set to `true` to show the trial LLM features. Defaults to hidden. |

## Production checklist

- Set `FARAID_ENV=production` and supply `FARAID_DATABASE_URL`.
- Run the API behind TLS (nginx/Caddy/load balancer).
- Restrict `/debug/vars` to internal networks only.
- Set `PUBLIC_LLM_ENABLED=true` only if you have wired in a real LLM
  provider by calling `server.WithLLM(...)` in the entrypoint.

## Schema migrations

Migrations are applied automatically by `faraidd` on startup using
`CREATE TABLE IF NOT EXISTS`. The migration source of truth is
`migrations/0001_create_cases.sql`. Manual rollback:

```sql
DROP TABLE cases;
```
