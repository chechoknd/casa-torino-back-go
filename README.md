# Casa Torino Backend

Backend en Go para Casa Torino, organizado con arquitectura hexagonal/limpia y preparado para conectarse a PostgreSQL existente.

## Requisitos

- Go 1.24+
- Docker y Docker Compose
- Acceso a una base de datos PostgreSQL

## Configuración

1. Copia `.env.example` a `.env`
2. Ajusta `DATABASE_URL` según tu entorno

## Arranque

```bash
make run
```

El servicio expone la API de clientes, productos, ingredientes, recetas, ordenes y pagos.

## Otros comandos

```bash
make down
make test
make sqlc
make migrate-up
make migrate-down
```

`DATABASE_URL` se usa para el backend dentro de Docker. `MIGRATIONS_DATABASE_URL` se usa desde tu host para `make migrate-up` y `make migrate-down`.

## Notas

- En esta fase no se crean ni modifican tablas existentes.
- `docker-compose.yml` incluye un contenedor PostgreSQL 15 para desarrollo local; si ya tienes una base de datos activa, sobrescribe `DATABASE_URL` en tu `.env`.
