# Database Connection Clarity

## Qué estaba pasando

El backend estaba escribiendo en la base de datos del contenedor Docker Compose, no en otra instancia PostgreSQL del host.

La URL real del backend es:

```text
postgres://user:password@db:5432/casa_torino?sslmode=disable
```

Eso significa:

- dentro de Docker, el backend usa el servicio `db`
- fuera de Docker, si quieres mirar esa misma base, debes conectarte al puerto publicado `5433`

## La misma base vista desde fuera de Docker

Usa esta URL:

```text
postgres://user:password@localhost:5433/casa_torino?sslmode=disable
```

## Credenciales de la base local de Docker

- host: `localhost`
- port: `5433`
- database: `casa_torino`
- user: `user`
- password: `password`

## Verificación rápida

### Ver conteos en la misma BD que usa Docker

```bash
make db-counts
```

### Abrir consola SQL en esa misma BD

```bash
make db-shell
```

## Archivo `.env`

Se dejó creado localmente con esta separación:

- `DATABASE_URL`
  Usa `db:5432` para el backend dentro de Docker
- `MIGRATIONS_DATABASE_URL`
  Usa `localhost:5433` para comandos lanzados desde tu host

## Conclusión

No hubo pérdida real de datos en la base del contenedor.

La confusión venía de consultar una instancia distinta de PostgreSQL o de no consultar la misma ruta que usa el backend.
