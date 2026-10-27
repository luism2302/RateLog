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
- `DB_PASSWORD`, `RATELOG_OWNER_PASSWORD` and `RATELOG_APP_PASSWORD`: three **different** passwords.

Protect the file and start the app:

```bash
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


### Updating

Change the image version in `compose.yaml` (or set `RATELOG_IMAGE` in `.env`) and run:

```bash
docker compose pull
docker compose up -d
```

New migrations are applied automatically before the app starts.

