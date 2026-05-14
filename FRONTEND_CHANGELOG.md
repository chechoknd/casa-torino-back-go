# Frontend Changelog - Customer Panel

Este archivo resume los cambios backend ya implementados que impactan al frontend. Se debe actualizar al finalizar cada feature o sprint para mantener claro que modulos del frontend deben ajustarse.

## Sprint 1 - Auth y Roles

**Estado:** Implementado

### Backend implementado

- Usuarios ahora tienen rol:
  - `ADMIN`
  - `CUSTOMER`
- La tabla `users` incluye columna `role`.
- Usuarios existentes fueron migrados con rol default `CUSTOMER`.
- `register` crea usuarios nuevos con rol `CUSTOMER`.
- `login` responde el rol del usuario.
- `refresh` mantiene el rol en la respuesta del usuario.
- El JWT ahora incluye `role`.
- Las rutas internas actuales requieren JWT y rol `ADMIN`.
- Si no hay token, backend responde `401 UNAUTHORIZED`.
- Si hay token valido pero el rol no tiene permiso, backend responde `403 FORBIDDEN`.

### Endpoints afectados

#### `POST /auth/register`

La respuesta del usuario ahora incluye `role`.

```json
{
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "demo",
    "full_name": "Demo User",
    "role": "CUSTOMER",
    "created_at": "2026-05-13T23:18:13Z"
  },
  "message": "ok"
}
```

#### `POST /auth/login`

La respuesta del login ahora incluye `role` dentro de `user`.

```json
{
  "data": {
    "access_token": "jwt",
    "refresh_token": "refresh-token",
    "token_type": "Bearer",
    "expires_at": "2026-05-13T23:33:17Z",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "username": "demo",
      "full_name": "Demo User",
      "role": "CUSTOMER",
      "created_at": "2026-05-13T23:18:13Z"
    }
  },
  "message": "ok"
}
```

#### `POST /auth/refresh`

La respuesta mantiene la misma estructura de `login` e incluye `user.role`.

### Reglas de permisos para frontend

- `ADMIN`:
  - Puede entrar al panel admin.
  - Puede consumir rutas internas actuales como `/products`, `/customers`, `/orders`, `/payments`, `/ingredients`, `/recipes`.

- `CUSTOMER`:
  - Puede autenticarse.
  - No debe entrar al panel admin.
  - Si intenta consumir rutas admin actuales recibira `403 FORBIDDEN`.
  - Su panel cliente queda para Sprint 2.

- `GUEST`:
  - Todavia no existe como usuario persistido.
  - Debe manejarse en frontend como modo anonimo sin JWT.
  - Rutas publicas de catalogo quedan para Sprint 2.

### Modulos frontend a modificar

- Auth service:
  - Guardar `user.role` despues de `login`, `register` y `refresh`.
  - Mantener envio de `Authorization: Bearer <token>` para rutas protegidas.

- Auth state/store:
  - Agregar campo `role` al usuario autenticado.
  - Derivar estado de sesion:
    - authenticated admin
    - authenticated customer
    - guest

- Guards de rutas:
  - Admin routes deben requerir `role === "ADMIN"`.
  - Customer routes deben requerir usuario autenticado con `role === "CUSTOMER"` o regla definida para admin.
  - Guest/catalog routes no deben requerir JWT.

- Interceptor HTTP:
  - `401 UNAUTHORIZED`: tratar como sesion no valida o no autenticada.
  - `403 FORBIDDEN`: tratar como usuario autenticado sin permisos.

- UI de login:
  - Despues de login redirigir segun rol:
    - `ADMIN` -> panel admin.
    - `CUSTOMER` -> panel cliente cuando exista.
  - Boton "Continuar como invitado" debe entrar sin token. Su funcionalidad real depende de Sprint 2.

### Pendiente para frontend

- Definir pantalla/ruta destino para `CUSTOMER`.
- Definir comportamiento visual de `403`.
- Definir como se creara/asignara el primer usuario `ADMIN`.
- Esperar Sprint 2 para consumir catalogo publico y guest mode real.

### Validacion realizada

- `GET /health` -> `200 OK`
- `POST /auth/register` -> responde usuario con `role: CUSTOMER`
- `POST /auth/login` -> responde JWT y usuario con `role`
- Ruta admin sin token -> `401 UNAUTHORIZED`
- Ruta admin con token `CUSTOMER` -> `403 FORBIDDEN`
- Ruta admin con token `ADMIN` -> `200 OK`

## Sprint 2 - Catalogo Publico y Guest Mode

**Estado:** Implementado

### Backend implementado

- Productos ahora soportan:
  - `image_url`
  - `is_public`
- Productos activos existentes fueron marcados como publicos para poblar el catalogo inicial.
- Productos inactivos no quedan publicos.
- Se agregaron rutas publicas para modo invitado.
- Las rutas publicas no requieren JWT.
- El catalogo publico no expone `cost_price`.

### Endpoints nuevos

#### `GET /public/products`

Lista productos activos y publicos.

Query params soportados:

- `product_type`: filtra por categoria/tipo existente.

Respuesta:

```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Almuerzo Ejecutivo Pollo",
      "description": "Arroz, pechuga a la plancha, ensalada y bebida",
      "product_type": "LUNCH",
      "base_price": "18500",
      "image_url": "",
      "created_at": "2026-05-07T00:20:49Z",
      "updated_at": "2026-05-07T00:20:49Z"
    }
  ],
  "message": "ok"
}
```

#### `GET /public/products/{id}`

Devuelve el detalle publico de un producto activo y visible.

No devuelve:

- `cost_price`
- `is_active`
- `is_public`

#### `GET /public/product-categories`

Lista categorias publicas derivadas de `product_type`.

Respuesta:

```json
{
  "data": ["CAKE", "EVENT", "JUICE", "LUNCH", "PLAN", "VACUUM_PACKED"],
  "message": "ok"
}
```

#### `GET /customer/profile`

Requiere JWT con rol `CUSTOMER`.

Devuelve el customer asociado por email al usuario autenticado.

Si el usuario existe pero no hay customer con el mismo email, responde `404 NOT_FOUND`.

#### `GET /customer/orders`

Requiere JWT con rol `CUSTOMER`.

Lista las ordenes del customer asociado por email al usuario autenticado.

Si no hay customer asociado, responde `404 NOT_FOUND`.

### Modulos frontend a modificar

- Catalog service:
  - Crear consumo de `GET /public/products`.
  - Crear consumo de `GET /public/products/{id}`.
  - Crear consumo de `GET /public/product-categories`.

- Guest mode:
  - El invitado debe navegar catalogo sin token.
  - No enviar `Authorization` en rutas `/public/*`.
  - No puede acceder a `/orders` ni `/customer/*`.

- Product models:
  - Agregar `image_url`.
  - Separar modelo admin de modelo catalogo publico.
  - El modelo publico no debe esperar `cost_price`.

- Product admin:
  - Agregar campos `image_url` e `is_public` al formulario de crear/editar producto.
  - Mostrar si un producto es publico o interno.

- Customer panel:
  - Consumir `GET /customer/profile` para datos del cliente.
  - Consumir `GET /customer/orders` para historial.
  - Requiere usuario autenticado con `role === "CUSTOMER"`.
  - Manejar `404 NOT_FOUND` como cuenta sin customer asociado.

### Pendiente para frontend

- Definir placeholder visual cuando `image_url` venga vacio.
- Definir filtros UI por `product_type`.
- Definir ruta/pantalla de detalle publico de producto.
- Definir pantalla para usuario customer sin customer asociado por email.
