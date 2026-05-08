# Autenticacion JWT

## Decisiones

- El backend usa JWT firmado con HMAC SHA-256 (`HS256`) y secreto provisto por `JWT_SECRET`.
- No se usa una dependencia externa para JWT; la implementacion usa la libreria estandar de Go para mantener baja la superficie de dependencias.
- Las contrasenas se almacenan solo como hash bcrypt. El costo se configura con `BCRYPT_COST`.
- Los endpoints publicos son `/auth/register` y `/auth/login`.
- Las rutas de negocio quedan protegidas cuando el router recibe un verificador JWT configurado. En tests de handlers se permite construir el router sin verificador para mantener pruebas unitarias simples.
- `REFRESH_TOKEN_EXPIRES` queda configurado para crecimiento futuro, pero esta implementacion entrega solo access token porque no hay flujo de refresh solicitado ni tabla de sesiones/tokens.

## Variables

```env
JWT_SECRET=secret-largo-y-aleatorio
JWT_EXPIRES_IN=15m
REFRESH_TOKEN_EXPIRES=168h
BCRYPT_COST=12
```

## Flujo

1. Registrar usuario con `POST /auth/register`.
2. Autenticar con `POST /auth/login`.
3. Usar `Authorization: Bearer <access_token>` para consumir rutas privadas.

## Instrucciones para el agente del frontend

Este backend ahora protege las rutas de negocio con JWT. El frontend debe implementar una capa de autenticacion antes de consumir endpoints como `/products`, `/customers`, `/orders`, `/payments`, `/ingredients` y `/recipes`.

### Contexto de despliegue

- Frontend Vercel: `https://casa-torino-front.vercel.app`
- Backend Render: usar la URL publica del servicio Render como `API_BASE_URL`.
- CORS ya permite:
  - `https://casa-torino-front.vercel.app`
  - `http://localhost:4200`

### Variables sugeridas en frontend

Usar la convencion del framework del frontend. Ejemplos:

```env
VITE_API_BASE_URL=https://<backend-render>.onrender.com
```

o, si usa Angular:

```env
NG_APP_API_BASE_URL=https://<backend-render>.onrender.com
```

No guardar `JWT_SECRET` en el frontend. Ese secreto es solo del backend.

### Contrato de API

Todas las respuestas exitosas vienen envueltas asi:

```json
{
  "data": {},
  "message": "ok"
}
```

Todas las respuestas de error vienen asi:

```json
{
  "error": "invalid credentials",
  "code": "INVALID_CREDENTIALS"
}
```

### Registro

Endpoint:

```http
POST /auth/register
Content-Type: application/json
```

Request:

```json
{
  "email": "admin@example.com",
  "username": "admin",
  "full_name": "Admin User",
  "password": "Password123"
}
```

Response `201 Created`:

```json
{
  "data": {
    "id": "6f218000-d49c-4eb5-bbb4-9417f7768345",
    "email": "admin@example.com",
    "username": "admin",
    "full_name": "Admin User",
    "created_at": "2026-05-08T05:16:04.069864Z"
  },
  "message": "ok"
}
```

Notas:

- El registro no devuelve token automaticamente.
- Despues de registrar, redirigir a login o ejecutar login con las mismas credenciales si el producto quiere iniciar sesion inmediatamente.

### Login

Endpoint:

```http
POST /auth/login
Content-Type: application/json
```

Request recomendado:

```json
{
  "email_or_username": "admin@example.com",
  "password": "Password123"
}
```

Tambien se acepta:

```json
{
  "email": "admin@example.com",
  "password": "Password123"
}
```

o:

```json
{
  "username": "admin",
  "password": "Password123"
}
```

Response `200 OK`:

```json
{
  "data": {
    "access_token": "<jwt>",
    "token_type": "Bearer",
    "expires_at": "2026-05-08T05:35:57.019264618Z",
    "user": {
      "id": "6f218000-d49c-4eb5-bbb4-9417f7768345",
      "email": "admin@example.com",
      "username": "admin",
      "full_name": "Admin User",
      "created_at": "2026-05-08T05:16:04.069864Z"
    }
  },
  "message": "ok"
}
```

### Consumo de rutas privadas

Enviar el token en todas las llamadas privadas:

```http
Authorization: Bearer <access_token>
```

Ejemplo:

```bash
curl https://<backend-render>.onrender.com/products \
  -H "Authorization: Bearer <access_token>"
```

Si falta el token o es invalido, el backend responde:

```http
401 Unauthorized
```

```json
{
  "error": "unauthorized",
  "code": "UNAUTHORIZED"
}
```

Si las credenciales son incorrectas en login:

```http
401 Unauthorized
```

```json
{
  "error": "invalid credentials",
  "code": "INVALID_CREDENTIALS"
}
```

Si el email ya existe en registro:

```http
409 Conflict
```

```json
{
  "error": "email already exists",
  "code": "DUPLICATE_EMAIL"
}
```

Si el username ya existe:

```http
409 Conflict
```

```json
{
  "error": "username already exists",
  "code": "DUPLICATE_USERNAME"
}
```

### Reglas de UI/UX esperadas

- Crear pantalla de login.
- Crear pantalla de registro si el producto lo necesita ahora; si no, dejar el flujo preparado.
- Guardar `access_token`, `expires_at` y `user` en un estado central de autenticacion.
- Agregar interceptor/client HTTP para inyectar `Authorization: Bearer <access_token>` en todas las rutas privadas.
- Si una respuesta privada devuelve `401`, limpiar sesion local y redirigir a login.
- No enviar token a dominios distintos del backend configurado.
- No mostrar mensajes tecnicos internos al usuario. Mapear errores:
  - `INVALID_CREDENTIALS`: credenciales invalidas.
  - `DUPLICATE_EMAIL`: el email ya esta registrado.
  - `DUPLICATE_USERNAME`: el usuario ya esta registrado.
  - `UNAUTHORIZED`: sesion expirada o no iniciada.

### Persistencia del token

Version pragmatica inicial:

- Persistir `access_token`, `expires_at` y `user` en `localStorage` o el mecanismo ya usado por el proyecto.
- Al iniciar la app, cargar la sesion si existe y si `expires_at` es futuro.
- Si `expires_at` ya paso, limpiar sesion y enviar a login.

Mejora futura:

- Migrar a cookies `HttpOnly` o refresh tokens si el backend agrega flujo de refresh/sesiones.

### Rutas protegidas en frontend

El frontend debe bloquear navegacion a pantallas privadas si no hay token valido. Como minimo:

- Productos
- Clientes
- Ingredientes
- Recetas
- Ordenes
- Pagos

La ruta de login debe estar disponible sin token.

### Checklist de implementacion frontend

1. Agregar variable `API_BASE_URL`.
2. Crear cliente HTTP centralizado.
3. Implementar `login(credentials)`.
4. Implementar `register(payload)` si aplica.
5. Crear store/context/service de auth con `token`, `expiresAt`, `user`, `login`, `logout`.
6. Inyectar header `Authorization` en requests privados.
7. Manejar `401` globalmente con logout y redirect a login.
8. Proteger rutas privadas en router.
9. Actualizar formularios para mostrar errores por `code`.
10. Probar CORS desde Vercel y desde `http://localhost:4200`.

### Comandos de prueba manual

Registro:

```bash
curl -X POST "$API_BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","username":"admin","full_name":"Admin User","password":"Password123"}'
```

Login:

```bash
curl -X POST "$API_BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email_or_username":"admin@example.com","password":"Password123"}'
```

Ruta privada:

```bash
curl "$API_BASE_URL/products" \
  -H "Authorization: Bearer <access_token>"
```

Preflight CORS esperado:

```bash
curl -i -X OPTIONS "$API_BASE_URL/products" \
  -H "Origin: https://casa-torino-front.vercel.app" \
  -H "Access-Control-Request-Method: GET" \
  -H "Access-Control-Request-Headers: Authorization, Content-Type"
```

Debe responder `204 No Content` con:

```http
Access-Control-Allow-Origin: https://casa-torino-front.vercel.app
Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization, X-Request-ID
```

### Consideraciones importantes

- El token actual es solo access token. No hay endpoint de refresh todavia.
- La expiracion se define en backend con `JWT_EXPIRES_IN`.
- Si el frontend conserva un token expirado, el backend respondera `401`.
- No asumir que `/products` u otras rutas son publicas: ahora requieren JWT.
- El backend responde siempre JSON, incluso en errores.
