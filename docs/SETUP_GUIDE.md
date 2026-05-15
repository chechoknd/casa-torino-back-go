# Casa Torino Backend — Guía de instalación y configuración

## Ramas activas

- `main` — producción
- `dev` — integración de features
- `feature/customer-panel-auth-roles` — feature actual (roles de autenticación para panel de cliente)

## Requisitos del sistema

- **Go** 1.24.2+
- **Docker** + **Docker Compose**
- **Git**
- **jq** (para scripts de seed)
- **Make**

## Clonar repositorio

```bash
git clone https://github.com/chechoknd/casa-torino-back-go.git
cd casa-torino-back-go
```

## Estructura del proyecto

```
.
├── cmd/api/main.go            # Entry point
├── internal/
│   ├── application/usecases   # Lógica de negocio (casos de uso)
│   ├── domain/                # Entidades e interfaces del dominio
│   ├── infrastructure/        # Config, DB, seguridad, SQLC generated
│   │   ├── config/            # Carga de variables de entorno
│   │   ├── database/
│   │   │   ├── postgres/      # Repositorios PostgreSQL (implementaciones)
│   │   │   └── sqlc/          # Código generado por SQLC
│   │   └── security/          # JWT, bcrypt
│   └── interfaces/http/       # Handlers, middleware, router
├── migrations/                # Migraciones SQL (golang-migrate)
├── queries/                   # Queries SQL para SQLC
├── scripts/                   # Scripts de seed
├── docker-compose.yml         # Servicios backend + PostgreSQL 15
├── Dockerfile                 # Multi-stage build (alpine)
├── Makefile                   # Comandos útiles
├── sqlc.yaml                  # Configuración de SQLC
├── go.mod / go.sum            # Dependencias
└── .env.example               # Template de variables de entorno
```

## Variables de entorno (.env)

Copia `.env.example` a `.env` y ajusta según tu entorno:

| Variable | Descripción | Ejemplo |
|---|---|---|
| `DATABASE_URL` | URL de conexión a PostgreSQL **dentro de Docker** | `postgres://user:password@db:5432/casa_torino?sslmode=disable` |
| `MIGRATIONS_DATABASE_URL` | URL de conexión **desde el host** para migraciones | `postgres://user:password@localhost:5433/casa_torino?sslmode=disable` |
| `FRONTEND_URL` | Origen permitido por CORS | `https://casa-torino-front.vercel.app` |
| `JWT_SECRET` | Clave secreta para firmar JWT (64+ chars recomendado) | `AQuEgzb2ovsg9IBiBZC9CDoBX2cxUMLFTJ3UD9oy11w=` |
| `JWT_EXPIRES_IN` | Tiempo de expiración del access token | `15m` |
| `REFRESH_TOKEN_EXPIRES` | Tiempo de expiración del refresh token | `168h` (7 días) |
| `BCRYPT_COST` | Costo de bcrypt para hash de contraseñas | `12` |
| `PORT` | Puerto donde escucha el servidor | `8080` |
| `ENV` | Entorno: `development` o `production` | `development` |

Nota: `http://localhost:4200` siempre está permitido por CORS además de `FRONTEND_URL`.

## Stack tecnológico

| Componente | Tecnología |
|---|---|
| Lenguaje | Go 1.24.2 |
| Router HTTP | go-chi/chi v5 |
| Base de datos | PostgreSQL 15 |
| Driver DB | jackc/pgx v5 |
| Migraciones | golang-migrate v4 |
| SQL a Go | sqlc v1.29 |
| JWT | golang-jwt v5 |
| Bcrypt | golang.org/x/crypto |
| Decimales | shopspring/decimal |
| Rate limiting | golang.org/x/time |
| Carga .env | joho/godotenv |
| Contenedor | Docker (multi-stage, alpine) |

## Dependencias Go

Ver `go.mod`. Todas las dependencias:

```
github.com/go-chi/chi/v5
github.com/golang-jwt/jwt/v5
github.com/google/uuid
github.com/jackc/pgx/v5
github.com/joho/godotenv
github.com/shopspring/decimal
golang.org/x/crypto
golang.org/x/time
```

## Instalación y arranque

### Opción 1: Docker Compose (recomendado)

```bash
# Construir y levantar backend + PostgreSQL
make run
# Equivalente a: docker compose up --build

# Detener
make down
# Equivalente a: docker compose down
```

Esto levanta:
- **backend**: compila y ejecuta el servidor en el puerto configurado (8080)
- **db**: PostgreSQL 15 en puerto 5433 (mapeado desde 5432 interno)

### Opción 2: Sin Docker

```bash
# 1. Asegúrate de tener PostgreSQL corriendo y configurar DATABASE_URL
# 2. Ejecutar migraciones (apunta a tu base de datos local)
make migrate-up

# 3. Descargar dependencias
go mod download

# 4. Compilar y ejecutar
go run ./cmd/api
# o compilar:
go build -o bin/backend ./cmd/api
./bin/backend
```

## Migraciones

```bash
# Aplicar migraciones
make migrate-up

# Revertir última migración
make migrate-down
```

Las migraciones se encuentran en `migrations/` usando `golang-migrate`. El orden actual:

1. `000001_create_core_tables` — tablas core (customers, products, ingredients, recipes, orders, etc.)
2. `000002_add_order_number`
3. `000003_create_users` — tabla de usuarios con email/password/username
4. `000004_create_refresh_tokens`
5. `000005_add_user_roles` — columna role en users (ADMIN, CUSTOMER)
6. `000006_add_product_catalog_fields` — image_url, is_public en products

## Generar código SQL

```bash
make sqlc
# Genera código Go en internal/infrastructure/database/sqlc/
# a partir de las queries en queries/ y el schema en migrations/
```

## Pruebas

```bash
make test
# Equivalente a: go test ./...
```

## Poblar datos de prueba

```bash
# Script bash simple (requiere API corriendo en localhost:8080 y jq instalado)
bash scripts/seed_data.sh

# Script mejorado con funciones helper
bash scripts/seed-api.sh
```

Ambos crean customers, productos, ingredientes, recetas, órdenes y pagos de ejemplo.

## Comandos útiles adicionales

```bash
# Ver conteo de registros por tabla
make db-counts

# Entrar a consola psql dentro del contenedor
make db-shell

# Verificar formato y estilo Go
go fmt ./...
go vet ./...
```

## Arquitectura de capas

```
cmd/api/main.go
  ↓
internal/interfaces/http/    ← Router, handlers, middleware
  ↓
internal/application/usecases  ← Casos de uso (orquestan lógica)
  ↓
internal/domain/             ← Entidades e interfaces (repositorios)
  ↓
internal/infrastructure/     ← Implementaciones concretas
  ├── database/postgres/     ← Repositorios PostgreSQL
  ├── database/sqlc/         ← Código generado por SQLC
  ├── config/                ← Carga de configuración (.env)
  └── security/              ← JWT y bcrypt
```

## API endpoints

### Públicos
| Método | Ruta | Descripción |
|---|---|---|
| `GET` | `/health` | Health check |
| `POST` | `/auth/register` | Registro de usuario |
| `POST` | `/auth/login` | Inicio de sesión |
| `POST` | `/auth/refresh` | Renovar access token |
| `POST` | `/auth/logout` | Cerrar sesión |
| `GET` | `/public/products` | Catálogo público de productos |
| `GET` | `/public/products/{id}` | Detalle de producto público |
| `GET` | `/public/product-categories` | Categorías de productos |

### Customer (requiere JWT con rol CUSTOMER)
| Método | Ruta | Descripción |
|---|---|---|
| `GET` | `/customer/profile` | Perfil del customer |
| `GET` | `/customer/orders` | Órdenes del customer |

### Admin (requiere JWT con rol ADMIN)
| Método | Ruta | Descripción |
|---|---|---|
| CRUD | `/customers` | Gestión de clientes |
| CRUD | `/products` | Gestión de productos |
| CRUD | `/ingredients` | Gestión de ingredientes |
| CRUD | `/recipes` | Gestión de recetas |
| CRUD | `/orders` | Gestión de órdenes |
| CRUD | `/payments` | Gestión de pagos |

## Regenerar desde cero (después de formatear PC)

```bash
# 1. Instalar Go 1.24+
# 2. Instalar Docker + Docker Compose
# 3. Clonar repo
git clone https://github.com/chechoknd/casa-torino-back-go.git
cd casa-torino-back-go

# 4. Configurar .env
cp .env.example .env
# Editar JWT_SECRET (generar uno nuevo o restaurar el anterior)

# 5. Levantar todo
make run

# 6. Aplicar migraciones (si no se ejecutan automáticamente)
make migrate-up

# 7. Opcional: sembrar datos de prueba
bash scripts/seed_data.sh

# 8. Verificar
curl localhost:8080/health
```
