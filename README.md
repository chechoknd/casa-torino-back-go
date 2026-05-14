# Casa Torino Backend

Backend en Go para Casa Torino, organizado con arquitectura hexagonal/limpia y preparado para conectarse a PostgreSQL existente.

## Requisitos

- Go 1.24+
- Docker y Docker Compose
- Acceso a una base de datos PostgreSQL

## Configuración

1. Copia `.env.example` a `.env`
2. Ajusta `DATABASE_URL`, `FRONTEND_URL` y `JWT_SECRET` según tu entorno

## Arranque

```bash
make run
```

El servicio expone la API de clientes, productos, ingredientes, recetas, ordenes, pagos, autenticacion, catalogo publico y panel cliente.

## Otros comandos

```bash
make down
make test
make sqlc
make migrate-up
make migrate-down
```

`DATABASE_URL` se usa para el backend dentro de Docker. `MIGRATIONS_DATABASE_URL` se usa desde tu host para `make migrate-up` y `make migrate-down`.
`FRONTEND_URL` define el origen permitido por CORS para el frontend desplegado; `http://localhost:4200` también queda permitido para desarrollo local.
`JWT_SECRET`, `JWT_EXPIRES_IN`, `REFRESH_TOKEN_EXPIRES` y `BCRYPT_COST` configuran autenticación y seguridad.

## Autenticación y roles

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","username":"admin","full_name":"Admin User","password":"Password123"}'

curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email_or_username":"admin@example.com","password":"Password123"}'
```

`/auth/login` tambien acepta `email` o `username` como alternativa a `email_or_username`.

Los usuarios tienen rol:

- `ADMIN`: puede consumir rutas internas de administracion.
- `CUSTOMER`: puede consumir rutas de panel cliente.

El modo invitado no crea usuario. Usa rutas publicas sin JWT.

Usa el token retornado en rutas privadas:

```bash
curl http://localhost:8080/products/ \
  -H "Authorization: Bearer <access_token>"
```

## Rutas principales

- Publicas:
  - `GET /health`
  - `POST /auth/register`
  - `POST /auth/login`
  - `POST /auth/refresh`
  - `POST /auth/logout`
  - `GET /public/products`
  - `GET /public/products/{id}`
  - `GET /public/product-categories`
- Customer:
  - `GET /customer/profile`
  - `GET /customer/orders`
- Admin:
  - `/customers`
  - `/products`
  - `/ingredients`
  - `/recipes`
  - `/orders`
  - `/payments`

## Notas

- Las migraciones agregan roles a `users` y campos de catalogo a `products`.
- `docker-compose.yml` incluye un contenedor PostgreSQL 15 para desarrollo local; si ya tienes una base de datos activa, sobrescribe `DATABASE_URL` en tu `.env`.
