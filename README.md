# RateLog

Search for video games, add them to your library, and write rated reviews.
Game data comes from the [RAWG API](https://rawg.io/apidocs).

## Installation with Docker

**Requirements:** Docker with the Compose plugin (Docker Desktop or Docker Engine) and a free RAWG API key.

```bash
git clone https://github.com/luism2302/RateLog.git
cd RateLog
cp .env.example .env
```

Fill in `.env`:

- `RAWG_API_KEY`: your RAWG API key.
- `DB_PASSWORD`, `RATELOG_OWNER_PASSWORD` and `RATELOG_APP_PASSWORD`: three **different** passwords. Generate each one
  with `openssl rand -hex 24`.

Protect the file and start the app:

```bash
chmod 600 .env
docker compose up -d
```

On the first run, Docker downloads the images and Postgres creates the database roles. Then the migrations are
applied and the app starts. Check the status with `docker compose ps`:

| Service   | Expected status                                         |
|-----------|---------------------------------------------------------|
| `db`      | `healthy`                                               |
| `migrate` | `exited (0)`: it applied the migrations and finished    |
| `app`     | `healthy`, a few seconds later                          |

Open **http://localhost:8000**.

> Session cookies are only sent over HTTPS, with `localhost` as the exception. This means that, without an HTTPS
> proxy in front, you can only sign in from the same machine, through `http://localhost:8000`.
>
> Don't change the port's `127.0.0.1` to `0.0.0.0` in `compose.yaml` to open it to other machines: the app would be
> exposed without HTTPS. To use it from elsewhere, put it behind a reverse proxy with HTTPS (for example, Caddy).

### Everyday use

```bash
docker compose ps              # status
docker compose logs -f app     # app logs
docker compose down            # stop (keeps your data)
docker compose down -v         # stop and DELETE all data
```

### Updating

Change the image version in `compose.yaml` (or set `RATELOG_IMAGE` in `.env`) and run:

```bash
docker compose pull
docker compose up -d
```

New migrations are applied automatically before the app starts.

### Troubleshooting

| Symptom | Fix |
|---|---|
| `... is missing from .env` | Fill in that variable in `.env`. |
| `migrate` fails with `password authentication failed` | The roles are only created the first time Postgres starts. If you changed the passwords afterwards, run `docker compose down -v` (**deletes your data**) and start again. |
| `db` doesn't start: `permission denied` on `01-roles.sh` | Make it executable with `chmod +x db/init/01-roles.sh`. Then run `docker compose down -v` and `docker compose up -d`. |
| Port 8000 is already in use | Set `RATELOG_PORT=8080` (or another port) in `.env`. |
| Search says the service is unavailable | Check `RAWG_API_KEY` in `.env` and run `docker compose up -d`. |

## Development

Requirements: Go, Node.js, [templ](https://templ.guide) and [goose](https://github.com/pressly/goose).

```bash
docker compose -f compose.dev.yaml up -d db   # Postgres on localhost:5433
make live                                     # templ + server + Tailwind in watch mode
```

In development, `.env` also needs `RATELOG_DB_DSN` and `RATELOG_MIGRATE_DSN` pointing to `localhost:5433`. Apply the
migrations with `go run ./cmd/app -migrate`.

To build the image: `docker build -t ratelog .`
