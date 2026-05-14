# Informe y manual de usuario API administrativa MVP

Fecha de revision: 2026-05-13

## Estado general

La aplicacion backend de Casa Torino esta implementada en Go con separacion por capas:

- `cmd/api`: arranque HTTP, configuracion y wiring de dependencias.
- `internal/interfaces/http`: rutas, handlers, requests, responses y middleware.
- `internal/application/usecases`: casos de uso de negocio.
- `internal/domain`: entidades, repositorios, errores y value objects.
- `internal/infrastructure`: configuracion, seguridad JWT/bcrypt y repositorios PostgreSQL.

El MVP administrativo cubre autenticacion y administracion de:

- Clientes.
- Productos.
- Ingredientes.
- Recetas.
- Ordenes.
- Pagos.

Estado de verificacion local:

- `go fmt ./...`: correcto.
- `go test ./...`: correcto.
- `go vet ./...`: correcto.
- `go build ./...`: correcto.

Nota: la validacion anterior confirma compilacion, pruebas unitarias/de handlers/de usecases y analisis estatico. No se ejecuto una prueba manual end-to-end contra un servidor HTTP vivo con base de datos sembrada en esta revision.

## Configuracion requerida

Variables de entorno principales:

| Variable | Uso |
| --- | --- |
| `DATABASE_URL` | Conexion PostgreSQL del backend. Obligatoria. |
| `FRONTEND_URL` | Origen permitido por CORS para frontend. |
| `JWT_SECRET` | Secreto para firmar/verificar JWT. Obligatoria. |
| `JWT_EXPIRES_IN` | Duracion del access token. Por defecto `15m`. |
| `REFRESH_TOKEN_EXPIRES` | Duracion esperada del refresh token. Por defecto `168h`. |
| `BCRYPT_COST` | Costo bcrypt. Por defecto `12` en Docker. |
| `PORT` | Puerto HTTP. Por defecto `8080`. |
| `ENV` | Ambiente de ejecucion. Por defecto `development`. |

Base URL local:

```text
http://localhost:8080
```

## Seguridad y formato comun

`GET /health` y `/auth/*` son publicos. El resto de endpoints administrativos requiere:

```text
Authorization: Bearer <access_token>
```

Formato de exito:

```json
{
  "data": {},
  "message": "ok"
}
```

Formato de error:

```json
{
  "error": "mensaje",
  "code": "ERROR_CODE"
}
```

Codigos de error principales:

| Code | HTTP | Significado |
| --- | ---: | --- |
| `INVALID_INPUT` | 400 | JSON invalido, UUID invalido, decimal invalido o datos requeridos incompletos. |
| `UNAUTHORIZED` | 401 | Token faltante/invalido o refresh token invalido. |
| `INVALID_CREDENTIALS` | 401 | Login incorrecto. |
| `NOT_FOUND` | 404 | Recurso no encontrado. |
| `INACTIVE` | 409 | Recurso existe pero esta desactivado/inactivo. |
| `DUPLICATE_EMAIL` | 409 | Email duplicado en cliente o usuario. |
| `DUPLICATE_USERNAME` | 409 | Username duplicado. |
| `INVALID_STATUS` | 422 | Estado o transicion de estado invalida. |
| `REQUEST_TOO_LARGE` | 413 | Body supera el limite configurado. |
| `INTERNAL_ERROR` | 500 | Error no clasificado. |

## Catalogos validos

### Customer type

- `PERSON`
- `COMPANY`

### Product type

- `LUNCH`
- `JUICE`
- `CAKE`
- `EVENT`
- `PLAN`
- `VACUUM_PACKED`

### Units

- `G`
- `ML`
- `UNIT`
- `KG`
- `L`

### Order status

- `PENDING`
- `CONFIRMED`
- `IN_PREPARATION`
- `READY`
- `DELIVERED`
- `CANCELLED`

Transiciones permitidas:

| Desde | Hacia |
| --- | --- |
| `PENDING` | `CONFIRMED`, `CANCELLED` |
| `CONFIRMED` | `IN_PREPARATION`, `CANCELLED` |
| `IN_PREPARATION` | `READY`, `CANCELLED` |
| `READY` | `DELIVERED` |
| `DELIVERED` | Sin transiciones nuevas |
| `CANCELLED` | Sin transiciones nuevas |

### Payment method

- `CASH`
- `TRANSFER`
- `NEQUI`
- `DAVIPLATA`
- `CARD`
- `OTHER`

### Payment status

- `PENDING`
- `PAID`
- `PARTIAL`
- `FAILED`
- `REFUNDED`

## Flujo de uso recomendado

1. Validar salud con `GET /health`.
2. Crear usuario administrativo con `POST /auth/register`.
3. Iniciar sesion con `POST /auth/login`.
4. Usar `access_token` en rutas administrativas.
5. Crear datos maestros: clientes, productos, ingredientes.
6. Crear recetas asociadas a productos y agregar ingredientes.
7. Crear ordenes, agregar items y avanzar estados.
8. Registrar pagos y consultar pagos por orden.

## Endpoints publicos

### Health

`GET /health`

Respuesta:

```json
{
  "status": "ok"
}
```

### Register

`POST /auth/register`

Body:

```json
{
  "email": "admin@example.com",
  "username": "admin",
  "full_name": "Admin User",
  "password": "Password123"
}
```

Reglas:

- `email`, `username`, `full_name` y `password` son obligatorios.
- `email` debe tener formato valido.
- `username` acepta letras, numeros, `_`, `.`, `-`, entre 3 y 50 caracteres.
- `password` debe tener minimo 8 caracteres.
- Email y username deben ser unicos.

Respuesta `201`:

```json
{
  "data": {
    "id": "uuid",
    "email": "admin@example.com",
    "username": "admin",
    "full_name": "Admin User",
    "created_at": "2026-05-13T00:00:00Z"
  },
  "message": "ok"
}
```

### Login

`POST /auth/login`

Body preferido:

```json
{
  "email_or_username": "admin@example.com",
  "password": "Password123"
}
```

Tambien acepta `email` o `username` como alternativa a `email_or_username`.

Respuesta `200`:

```json
{
  "data": {
    "access_token": "jwt",
    "refresh_token": "token",
    "token_type": "Bearer",
    "expires_at": "2026-05-13T00:15:00Z",
    "user": {
      "id": "uuid",
      "email": "admin@example.com",
      "username": "admin",
      "full_name": "Admin User",
      "created_at": "2026-05-13T00:00:00Z"
    }
  },
  "message": "ok"
}
```

### Refresh token

`POST /auth/refresh`

Body:

```json
{
  "refresh_token": "token"
}
```

Respuesta: nuevo `access_token` y nuevo `refresh_token`. El refresh token anterior se revoca.

### Logout

`POST /auth/logout`

Body:

```json
{
  "refresh_token": "token"
}
```

Respuesta:

```json
{
  "data": {
    "status": "ok"
  },
  "message": "ok"
}
```

## Clientes

Endpoints:

| Metodo | Ruta | Uso |
| --- | --- | --- |
| `POST` | `/customers/` | Crear cliente. |
| `GET` | `/customers/` | Listar clientes activos. |
| `GET` | `/customers/{id}` | Obtener cliente activo. |
| `PUT` | `/customers/{id}` | Actualizar cliente. |
| `DELETE` | `/customers/{id}` | Desactivar cliente. |

Body crear/actualizar:

```json
{
  "full_name": "Cliente Demo",
  "phone": "3001234567",
  "email": "cliente.demo@example.com",
  "customer_type": "PERSON"
}
```

Propiedades de respuesta:

| Campo | Tipo | Descripcion |
| --- | --- | --- |
| `id` | UUID | Identificador del cliente. |
| `full_name` | string | Nombre completo o razon social. |
| `phone` | string | Telefono. |
| `email` | string | Email normalizado a minusculas. |
| `customer_type` | enum | `PERSON` o `COMPANY`. |
| `created_at` | datetime | Fecha de creacion. |
| `updated_at` | datetime | Ultima actualizacion. |
| `is_active` | boolean | Estado logico. |

Reglas:

- `full_name`, `phone`, `email` y `customer_type` son obligatorios.
- El email no puede duplicarse.
- `DELETE` no borra fisicamente: desactiva el registro.

## Productos

Endpoints:

| Metodo | Ruta | Uso |
| --- | --- | --- |
| `POST` | `/products/` | Crear producto. |
| `GET` | `/products/` | Listar productos activos. |
| `GET` | `/products/?product_type=LUNCH` | Listar por tipo. |
| `GET` | `/products/{id}` | Obtener producto activo. |
| `PUT` | `/products/{id}` | Actualizar producto. |
| `DELETE` | `/products/{id}` | Desactivar producto. |

Body crear/actualizar:

```json
{
  "name": "Almuerzo Ejecutivo",
  "description": "Arroz, proteina y ensalada",
  "product_type": "LUNCH",
  "base_price": "18000",
  "cost_price": "9000"
}
```

Propiedades de respuesta:

| Campo | Tipo | Descripcion |
| --- | --- | --- |
| `id` | UUID | Identificador del producto. |
| `name` | string | Nombre comercial. |
| `description` | string | Descripcion. |
| `product_type` | enum | Tipo de producto. |
| `base_price` | decimal | Precio de venta. |
| `cost_price` | decimal | Costo del producto. |
| `is_active` | boolean | Estado logico. |
| `created_at` | datetime | Fecha de creacion. |
| `updated_at` | datetime | Ultima actualizacion. |

Reglas:

- `name` es obligatorio.
- `base_price` debe ser positivo.
- `base_price` y `cost_price` se envian como string decimal.
- El filtro `product_type` debe usar un valor valido.
- `DELETE` desactiva el producto.

## Ingredientes

Endpoints:

| Metodo | Ruta | Uso |
| --- | --- | --- |
| `POST` | `/ingredients/` | Crear ingrediente. |
| `GET` | `/ingredients/` | Listar ingredientes activos. |
| `GET` | `/ingredients/{id}` | Obtener ingrediente activo. |
| `PUT` | `/ingredients/{id}` | Actualizar ingrediente. |
| `DELETE` | `/ingredients/{id}` | Desactivar ingrediente. |

Body crear/actualizar:

```json
{
  "name": "Arroz",
  "unit": "KG",
  "average_cost": "4500",
  "stock": "20",
  "minimum_stock": "5"
}
```

Propiedades de respuesta:

| Campo | Tipo | Descripcion |
| --- | --- | --- |
| `id` | UUID | Identificador del ingrediente. |
| `name` | string | Nombre. |
| `unit` | enum | Unidad base. |
| `average_cost` | decimal | Costo promedio. |
| `stock` | decimal | Existencia actual. |
| `minimum_stock` | decimal | Stock minimo esperado. |
| `created_at` | datetime | Fecha de creacion. |
| `updated_at` | datetime | Ultima actualizacion. |
| `is_active` | boolean | Estado logico. |

Reglas:

- `name` es obligatorio.
- `unit` debe estar dentro del catalogo de unidades.
- Valores numericos se envian como string decimal.
- `DELETE` desactiva el ingrediente.

## Recetas

Endpoints:

| Metodo | Ruta | Uso |
| --- | --- | --- |
| `POST` | `/recipes/` | Crear receta para un producto. |
| `GET` | `/recipes/` | Listar recetas. |
| `POST` | `/recipes/{id}/items` | Agregar ingrediente a receta. |
| `GET` | `/recipes/product/{product_id}` | Obtener receta por producto. |
| `GET` | `/recipes/{id}/cost` | Calcular costo de receta. |

Body crear receta:

```json
{
  "product_id": "uuid",
  "name": "Receta Almuerzo Ejecutivo",
  "portions": 1
}
```

Body agregar item:

```json
{
  "ingredient_id": "uuid",
  "quantity": "0.25",
  "unit": "KG"
}
```

Propiedades de respuesta de receta:

| Campo | Tipo | Descripcion |
| --- | --- | --- |
| `id` | UUID | Identificador de receta. |
| `product_id` | UUID | Producto asociado. |
| `product_name` | string | Nombre del producto asociado. |
| `name` | string | Nombre de receta. |
| `portions` | integer | Porciones que produce. |
| `items` | array | Ingredientes de la receta. |
| `created_at` | datetime | Fecha de creacion. |
| `updated_at` | datetime | Ultima actualizacion. |
| `is_active` | boolean | Estado logico. |

Propiedades de `items`:

| Campo | Tipo | Descripcion |
| --- | --- | --- |
| `id` | UUID | Identificador del item. |
| `recipe_id` | UUID | Receta asociada. |
| `ingredient_id` | UUID | Ingrediente asociado. |
| `quantity` | decimal | Cantidad usada. |
| `unit` | enum | Unidad usada. |

Respuesta de costo:

```json
{
  "data": {
    "recipe_id": "uuid",
    "cost": "12345.67"
  },
  "message": "ok"
}
```

Reglas:

- La receta requiere producto activo.
- `name` es obligatorio.
- `portions` debe ser mayor a cero.
- Para agregar item, receta e ingrediente deben estar activos.
- `quantity` debe ser mayor a cero.
- `unit` debe estar dentro del catalogo valido.

## Ordenes

Endpoints:

| Metodo | Ruta | Uso |
| --- | --- | --- |
| `POST` | `/orders/` | Crear orden. |
| `GET` | `/orders/` | Listar ordenes. |
| `GET` | `/orders/?customer_id={uuid}` | Listar ordenes por cliente. |
| `GET` | `/orders/{id}` | Obtener orden. |
| `POST` | `/orders/{id}/items` | Agregar producto a orden. |
| `PATCH` | `/orders/{id}/status` | Cambiar estado de orden. |
| `GET` | `/orders/{id}/payments` | Consultar pagos de la orden. |

Body crear orden:

```json
{
  "customer_id": "uuid",
  "discount": "1000"
}
```

Body agregar item:

```json
{
  "product_id": "uuid",
  "quantity": 2
}
```

Body actualizar estado:

```json
{
  "status": "CONFIRMED"
}
```

Propiedades de respuesta:

| Campo | Tipo | Descripcion |
| --- | --- | --- |
| `id` | UUID | Identificador de orden. |
| `customer_id` | UUID | Cliente asociado. |
| `customer_name` | string | Nombre del cliente. |
| `order_number` | integer | Consecutivo generado por base de datos. |
| `order_label` | string | Consecutivo formateado, ejemplo `#0001`. |
| `status` | enum | Estado de orden. |
| `items` | array | Productos de la orden. |
| `subtotal` | decimal | Suma de items. |
| `discount` | decimal | Descuento aplicado. |
| `total` | decimal | Total final. |
| `created_at` | datetime | Fecha de creacion. |
| `updated_at` | datetime | Ultima actualizacion. |

Propiedades de `items`:

| Campo | Tipo | Descripcion |
| --- | --- | --- |
| `id` | UUID | Identificador del item. |
| `order_id` | UUID | Orden asociada. |
| `product_id` | UUID | Producto asociado. |
| `product_name` | string | Nombre del producto. |
| `quantity` | integer | Cantidad. |
| `unit_price` | decimal | Precio unitario tomado del producto. |
| `subtotal` | decimal | `unit_price * quantity`. |

Reglas:

- La orden requiere cliente activo.
- `discount` no puede ser negativo.
- La orden nueva inicia en `PENDING`.
- El producto agregado debe estar activo.
- `quantity` debe ser mayor a cero.
- Las transiciones de estado deben respetar el flujo permitido.

## Pagos

Endpoints:

| Metodo | Ruta | Uso |
| --- | --- | --- |
| `GET` | `/payments/` | Listar pagos. |
| `POST` | `/payments/` | Crear pago. |
| `PATCH` | `/payments/{id}/status` | Actualizar estado de pago. |
| `GET` | `/orders/{id}/payments` | Listar pagos por orden. |

Body crear pago:

```json
{
  "order_id": "uuid",
  "amount": "35000",
  "method": "CASH",
  "status": "PENDING"
}
```

Body actualizar estado:

```json
{
  "status": "PAID"
}
```

Propiedades de respuesta:

| Campo | Tipo | Descripcion |
| --- | --- | --- |
| `id` | UUID | Identificador del pago. |
| `order_id` | UUID | Orden asociada. |
| `order_number` | integer | Consecutivo de orden. |
| `order_label` | string | Consecutivo formateado. |
| `amount` | decimal | Valor pagado. |
| `method` | enum | Metodo de pago. |
| `status` | enum | Estado del pago. |
| `products` | array | Productos incluidos en la orden pagada. |
| `created_at` | datetime | Fecha de creacion. |
| `updated_at` | datetime | Ultima actualizacion. |

Propiedades de `products`:

| Campo | Tipo | Descripcion |
| --- | --- | --- |
| `product_id` | UUID | Producto asociado. |
| `product_name` | string | Nombre del producto. |
| `quantity` | integer | Cantidad vendida en la orden. |

Reglas:

- El pago requiere orden existente.
- `amount` debe ser mayor a cero.
- `method` debe estar dentro del catalogo de metodos de pago.
- `status` debe estar dentro del catalogo de estados de pago.

## Estado funcional por modulo

| Modulo | Estado | Evidencia |
| --- | --- | --- |
| Salud | Operativo | Ruta directa `GET /health`. |
| Autenticacion | Operativo | Handlers, usecase y pruebas pasan. Incluye register, login, refresh y logout. |
| Clientes | Operativo | CRUD logico con validaciones y pruebas de handler/usecase. |
| Productos | Operativo | CRUD logico, filtro por tipo y pruebas de handler/usecase. |
| Ingredientes | Operativo | CRUD logico y pruebas de handler/usecase. |
| Recetas | Operativo | Crear, listar, agregar items, buscar por producto y calcular costo. |
| Ordenes | Operativo | Crear, listar, filtrar por cliente, agregar items y cambiar estado. |
| Pagos | Operativo | Crear, listar, listar por orden y actualizar estado. |

## Observaciones para cierre de MVP administrativo

- Los endpoints administrativos estan protegidos por JWT. Cualquier documentacion, script o coleccion Postman debe incluir `Authorization: Bearer <access_token>`.
- `scripts/seed-api.sh` no agrega header de autenticacion, por lo que requiere ajuste antes de usarse contra la API actual protegida.
- Los decimales se reciben como strings en requests para evitar problemas de precision.
- Los deletes son desactivaciones logicas para clientes, productos e ingredientes.
- No hay endpoints de actualizacion/eliminacion para recetas ni eliminacion de items de receta/orden en el router actual.
- No hay endpoint publico de carrito de compras todavia. La siguiente fase puede apoyarse en productos, ordenes e items de orden, pero probablemente necesitara una capa separada de carrito/sesion antes de convertir a orden.
- La integracion con WhatsApp no existe aun en el backend. Para esa fase conviene definir primero eventos de negocio: orden creada, orden confirmada, pago recibido y cambios de estado.

## Checklist sugerido antes de avanzar a carrito y WhatsApp

- Actualizar coleccion Postman/scripts con autenticacion.
- Ejecutar smoke test HTTP end-to-end contra Docker con base migrada y usuario admin.
- Definir contrato del carrito: anonimo vs autenticado, persistencia, expiracion y conversion a orden.
- Definir estrategia WhatsApp: proveedor, plantillas, eventos disparadores, reintentos y auditoria.
- Agregar tests de integracion para el flujo principal: login, crear cliente/producto, crear orden, agregar item, registrar pago.
