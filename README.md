# RateLog

Busca videojuegos, guárdalos en tu biblioteca y escribe reseñas con calificación.
Los datos de los juegos vienen de la [API de RAWG](https://rawg.io/apidocs).

## Instalación con Docker

**Requisitos:** Docker con el plugin Compose (Docker Desktop o Docker Engine) y una clave gratuita de la API de RAWG.

```bash
git clone https://github.com/luism2302/RateLog.git
cd RateLog
cp .env.example .env
```

Rellena `.env`:

- `RAWG_API_KEY`: tu clave de RAWG.
- `DB_PASSWORD`, `RATELOG_OWNER_PASSWORD` y `RATELOG_APP_PASSWORD`: tres contraseñas **distintas**. Genera cada una con
  `openssl rand -hex 24`.

Protege el archivo y arranca:

```bash
chmod 600 .env
docker compose up -d
```

La primera vez, Docker descarga las imágenes y Postgres crea los roles de la base de datos. Después se aplican las
migraciones y arranca la app. Comprueba el estado con `docker compose ps`:

| Servicio  | Estado esperado                                         |
|-----------|---------------------------------------------------------|
| `db`      | `healthy`                                               |
| `migrate` | `exited (0)`: aplicó las migraciones y terminó          |
| `app`     | `healthy`, unos segundos después                        |

Abre **http://localhost:8000**.

> Las cookies de sesión solo viajan por HTTPS, con la excepción de `localhost`. Por eso, sin un proxy con HTTPS
> delante, solo podrás iniciar sesión desde la misma máquina, entrando por `http://localhost:8000`.
>
> No cambies el `127.0.0.1` del puerto en `compose.yaml` por `0.0.0.0` para abrirla a otras máquinas: la app quedaría
> expuesta sin HTTPS. Para usarla desde fuera, ponla detrás de un proxy inverso con HTTPS (por ejemplo, Caddy).

### Uso diario

```bash
docker compose ps              # estado
docker compose logs -f app     # logs de la app
docker compose down            # parar (conserva los datos)
docker compose down -v         # parar y BORRAR todos los datos
```

### Actualizar

Cambia la versión de la imagen en `compose.yaml` (o define `RATELOG_IMAGE` en `.env`) y ejecuta:

```bash
docker compose pull
docker compose up -d
```

Las migraciones nuevas se aplican solas antes de que arranque la app.

### Problemas frecuentes

| Síntoma | Solución |
|---|---|
| `falta ... en .env` | Rellena esa variable en `.env`. |
| `migrate` falla con `password authentication failed` | Los roles se crean solo la primera vez que arranca Postgres. Si cambiaste las contraseñas después, ejecuta `docker compose down -v` (**borra los datos**) y vuelve a arrancar. |
| `db` no arranca: `permission denied` en `01-roles.sh` | Dale permiso de ejecución con `chmod +x db/init/01-roles.sh`. Después, `docker compose down -v` y `docker compose up -d`. |
| El puerto 8000 está ocupado | Pon `RATELOG_PORT=8080` (u otro puerto) en `.env`. |
| La búsqueda dice que el servicio no está disponible | Revisa `RAWG_API_KEY` en `.env` y ejecuta `docker compose up -d`. |

## Desarrollo

Requisitos: Go, Node.js, [templ](https://templ.guide) y [goose](https://github.com/pressly/goose).

```bash
docker compose -f compose.dev.yaml up -d db   # Postgres en localhost:5433
make live                                     # templ + servidor + Tailwind en modo watch
```

En desarrollo, `.env` necesita además `RATELOG_DB_DSN` y `RATELOG_MIGRATE_DSN` apuntando a `localhost:5433`. Las
migraciones se aplican con `go run ./cmd/app -migrate`.

Para construir la imagen: `docker build -t ratelog .`
