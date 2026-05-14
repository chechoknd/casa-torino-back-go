# Vulnerabilities & Risks Analysis

> **Fecha:** 2026-05-08  
> **Proyecto:** Casa Torino Backend (Go)  
> **Objetivo:** Identificar vulnerabilidades, riesgos de seguridad y brechas de hardening antes de implementar protecciones.

---

## Resumen ejecutivo (cambios detectados)

Desde el análisis inicial se identificaron **múltiples mejoras de seguridad ya implementadas** en el código.  
A continuación el detalle de vulnerabilidades **corregidas** y las que **persisten**.

| Prioridad | Vulnerabilidad | Estado |
|---|---|---|
| 🔴 Crítica | Sin rate limiting en auth | ✅ **CORREGIDO** |
| 🔴 Crítica | Sin límite de body HTTP | ✅ **CORREGIDO** |
| 🔴 Crítica | `JWT_SECRET` placeholder en `.env.example` | ⚠️ **PERSISTE** (documentación, no código) |
| 🟠 Alta | Sin `WriteTimeout` / `IdleTimeout` | ✅ **CORREGIDO** |
| 🟠 Alta | Contenedor Docker como root | ✅ **CORREGIDO** |
| 🟠 Alta | Headers de seguridad HTTP faltantes | ✅ **CORREGIDO** |
| 🟠 Alta | `sslmode=disable` en DB | ⚠️ **PERSISTE** (local dev) |
| 🟡 Media | JWT custom (vs librería estándar) | ✅ **CORREGIDO** (migrado a `golang-jwt/jwt/v5`) |
| 🟡 Media | Sin refresh token rotation | ✅ **CORREGIDO** |
| 🟡 Media | Sin endpoint healthcheck | ✅ **CORREGIDO** |
| 🟡 Media | Logger no estructurado / sin audit log | ⚠️ **PERSISTE** |
| 🟢 Baja | Sin `nbf`, `jti` en JWT | ⚠️ **PERSISTE** (bajo riesgo) |
| 🟢 Baja | Sin MFA | ⚠️ **PERSISTE** (depende del contexto) |

**Resumen:** 8 de 13 vulnerabilidades fueron corregidas. Quedan 5 pendientes (2 de riesgo medio-bajo aceptables, 1 propia de entorno local, 1 de documentación y 1 de logging).

---

## Estado del proyecto

| Componente | Estado |
|---|---|
| Arquitectura | Hexagonal / Clean Architecture |
| Router | `go-chi/chi/v5` |
| DB Driver | `pgx/v5` (connection pool) |
| ORM / Codegen | `sqlc` (queries parametrizadas) |
| Auth | JWT custom (HMAC-SHA256) + Bcrypt |
| Tests | Handler + UseCase + JWT + CORS + Auth middleware |
| Docker | Multi-stage build + docker-compose dev |
| CI/CD | No detectado |

---

## Cambios detectados vs análisis inicial

### Vulnerabilidades corregidas

| # | Vulnerabilidad original | Cambio detectado | Archivos |
|---|---|---|---|
| 1 | **Sin rate limiting** | Se implementó `RateLimiter` middleware con límite default 100 req/min y 5 req/min en `/auth/` | `middleware/rate_limiter.go`, `routes/router.go:32-35` |
| 2 | **Sin límite de body HTTP** | Se implementó `MaxBodyBytes` middleware con límite de 1MB (`1 << 20`) | `middleware/max_body_bytes.go`, `routes/router.go:31` |
| 3 | **Sin `WriteTimeout` / `IdleTimeout`** | Configurados: `ReadTimeout: 15s`, `WriteTimeout: 15s`, `IdleTimeout: 60s` | `cmd/api/main.go:97-100` |
| 4 | **Headers de seguridad faltantes** | Se implementó `SecurityHeaders` middleware: `X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`, `HSTS` (condicional a HTTPS), `Referrer-Policy` | `middleware/security_headers.go`, `routes/router.go:30` |
| 5 | **Contenedor Docker como root** | Se agregó usuario no-root: `adduser -D -u 1001 appuser` + `USER appuser` | `Dockerfile:22-23` |
| 6 | **JWT custom** | Migrado de implementación manual (`crypto/hmac`) a `github.com/golang-jwt/jwt/v5` | `internal/infrastructure/security/jwt.go` (reescrito completo), `go.mod` |
| 7 | **Sin refresh token rotation** | Implementado: login emite refresh token, `/auth/refresh` rota (revoca + emite nuevo), `/auth/logout` revoca | `auth/auth.go:147-218`, `entities/refresh_token.go`, `repositories/refresh_token_repository.go`, `postgres/refresh_token_repository.go` |
| 8 | **Sin endpoint healthcheck** | Agregado `GET /health` con respuesta `{"status":"ok"}` | `routes/router.go:40-43` |

### Nuevos endpoints de auth detectados

| Ruta | Método | Función |
|---|---|---|
| `/auth/refresh` | POST | Rota refresh token (revoca anterior, emite nuevo access + refresh) |
| `/auth/logout` | POST | Revoca el refresh token activo |

### Vulnerabilidades que persisten

| # | Vulnerabilidad | Razón | Recomendación |
|---|---|---|---|
| 1 | `JWT_SECRET=change-this-secret-in-production` en `.env.example` | Es un placeholder documentado, no código ejecutable. Riesgo solo si alguien despliega sin cambiarlo | Agregar validación en `config.go` que rechace el placeholder |
| 2 | `sslmode=disable` en DB URL | Solo afecta entorno local. Producción debe usar `sslmode=require` | Documentar en README |
| 3 | Logger usa `log.Printf` sin estructura | No impide funcionamiento ni abre vulnerabilidad | Migrar a `slog` o `zerolog` cuando se haga refactor de observabilidad |
| 4 | Sin `nbf` (not before) / `jti` (token ID) en JWT | Bajo impacto. La librería `golang-jwt/jwt/v5` soporta `RegisteredClaims` que incluye `ID` | Agregar si se requiere revocación individual de tokens |
| 5 | Sin MFA | Depende del contexto de negocio. No es blocker para release inicial | Evaluar según requerimientos |

---

## 1. Credenciales y secretos expuestos

### 1.1 `.env` con credenciales hardcodeadas

| Archivo | Línea | Contenido | Riesgo |
|---|---|---|---|
| `.env` | 1 | `DATABASE_URL=postgres://user:password@db:5432/casa_torino?sslmode=disable` | Alto |
| `.env` | 2 | `MIGRATIONS_DATABASE_URL=postgres://user:password@localhost:5433/casa_torino?sslmode=disable` | Alto |
| `docker-compose.yml` | 25-26 | `POSTGRES_USER: user`, `POSTGRES_PASSWORD: password` | Medio |

> **Nota:** `.env` NO está trackeado por git. Solo aplica a entornos donde alguien copie el repo y ejecute sin cambiar credenciales.

### 1.2 `.env.example` con placeholder peligroso

| Línea | Contenido | Riesgo |
|---|---|---|
| 4 | `JWT_SECRET=change-this-secret-in-production` | **Crítico si se despliega sin cambiar** |

---

## 2. Configuración del servidor HTTP

| Hallazgo | Archivo | Riesgo | Estado |
|---|---|---|---|
| `WriteTimeout: 15s` | `cmd/api/main.go:99` | - | ✅ Configurado |
| `IdleTimeout: 60s` | `cmd/api/main.go:100` | - | ✅ Configurado |
| `ReadTimeout: 15s` | `cmd/api/main.go:97` | - | ✅ Configurado |
| `ReadHeaderTimeout: 5s` | `cmd/api/main.go:98` | - | ✅ Configurado |
| Límite de body: 1MB (`1 << 20`) | `middleware/max_body_bytes.go` + `router.go:31` | - | ✅ Implementado |
| Rate limiting: 100 req/min (default), 5 req/min (`/auth/`) | `middleware/rate_limiter.go` + `router.go:32-35` | - | ✅ Implementado |
| Sin TLS/HTTPS | Global | Alto (requiere proxy reverso) | ⚠️ **PERSISTE** |

> Los 6 hallazgos de riesgo del análisis inicial han sido corregidos. Solo resta TLS, que es responsabilidad del proxy reverso en producción.

---

## 3. Seguridad en JWT

| Hallazgo | Archivo | Riesgo | Estado |
|---|---|---|---|
| Librería estándar `golang-jwt/jwt/v5` | `internal/infrastructure/security/jwt.go` | - | ✅ Migrado |
| Algoritmo validado (solo HS256) | `jwt.go:69` | - | ✅ Bueno |
| Expiración requerida (`WithExpirationRequired`) | `jwt.go:67` | - | ✅ Bueno |
| `WithValidMethods` previene algorithm confusion | `jwt.go:68` | - | ✅ Bueno |
| **Sin `nbf` (not before)** | - | Bajo | ⚠️ **PERSISTE** |
| **Sin `jti` (token ID)** | - | Bajo | ⚠️ **PERSISTE** |
| Refresh token rotation implementado | `auth/auth.go:147-218` | - | ✅ Implementado |

> La migración a `golang-jwt/jwt/v5` elimina el riesgo de errores en implementación manual. La librería es la más usada en el ecosistema Go.

---

## 4. Headers de seguridad HTTP

| Header | Valor | Estado |
|---|---|---|
| `X-Frame-Options` | `DENY` | ✅ Implementado |
| `X-Content-Type-Options` | `nosniff` | ✅ Implementado |
| `Strict-Transport-Security` | `max-age=63072000; includeSubDomains` (solo si HTTPS) | ✅ Implementado |
| `Referrer-Policy` | `strict-origin-when-cross-origin` | ✅ Implementado |
| `Cache-Control` | `no-store` (vía helper) | ✅ Implementado |
| `Content-Type` | `application/json` (middleware) | ✅ Implementado |
| `X-XSS-Protection` | - | Bajo (obsoleto en navegadores modernos) |
| `Content-Security-Policy` | - | Bajo (API JSON, no HTML) |

> 5 de 5 headers prioritarios fueron implementados vía middleware `SecurityHeaders` + middleware `ContentType` + helper `noCache()`.

---

## 5. Docker y despliegue

| Hallazgo | Archivo | Riesgo | Estado |
|---|---|---|---|
| Usuario no-root (`appuser`, UID 1001) | `Dockerfile:22-23` | - | ✅ Implementado |
| `sslmode=disable` en DB_URL | `.env`, `docker-compose.yml` | Alto | ⚠️ **PERSISTE** (solo local) |
| Sin `HEALTHCHECK` | `Dockerfile` | Bajo | ⚠️ **PERSISTE** |
| Sin secrets management | Global | Alto | ⚠️ **PERSISTE** (requiere infraestructura) |

> El cambio a usuario no-root elimina el riesgo de escalación de privilegios desde el contenedor.

---

## 6. Autenticación y sesiones

| Hallazgo | Riesgo |
|---|---|
| Sin rate limiting en login | Alto |
| Sin bloqueo por intentos fallidos | Alto |
| Sin refresh tokens implementados | Medio |
| Sin registro de intentos de login fallidos (audit log) | Medio |
| Sin MFA / 2FA | Bajo (depende del contexto) |

### Riesgo concreto

Un atacante puede hacer fuerza bruta contra `/auth/login` sin restricción. La única defensa es la complejidad de la contraseña y el coste de bcrypt (configurable, default 10).

---

## 7. Validación de entrada

| Hallazgo | Archivo | Riesgo |
|---|---|---|
| Email validado con `mail.ParseAddress` | `auth/auth.go:141` | Bueno |
| Username validado con regex `^[a-zA-Z0-9_.-]{3,50}$` | `auth/auth.go:21` | Bueno |
| Password mínimo 8 caracteres | `auth/auth.go:19` | Bueno |
| Normalización: lowercase + trim | `auth/auth.go:154-159` | Bueno |
| Validación básica en customer/product/ingredient | usecases/ | Bueno |
| **Sin límite de tamaño en campos de texto** | general | Medio |
| **Sin sanitización de input para XSS en campos de texto** | general | Bajo (API JSON) |

---

## 8. SQL Injection

**No se detectó riesgo.** Todas las consultas usan placeholders `$1`, `$2` (sqlc + raw pgx). No hay concatenación de strings en SQL.

---

## 9. CORS

| Hallazgo | Riesgo |
|---|---|
| Whitelist de orígenes (`FrontendURL` + `localhost:4200`) | Bueno |
| Sin `Access-Control-Allow-Credentials: true` | Bueno (no expone cookies) |
| Métodos y headers acotados | Bueno |

---

## 10. Logger y monitoreo

| Hallazgo | Riesgo |
|---|---|
| Logger usa `log.Printf` en vez de logger estructurado | Bajo |
| No se registran intentos de auth fallidos | Medio |
| No hay endpoint `/health` o `/ready` | Medio |
| No hay métricas de seguridad | Medio |

---

## 11. Base de datos

| Hallazgo | Riesgo |
|---|---|
| `sslmode=disable` en conexión | Alto |
| Pool de conexiones configurado (pgxpool) | Bueno |
| Transacciones en order/recipe repositories | Bueno |
| Mapeo de errores de BD (código 23505, 23503) | Bueno |
| Sin migraciones automáticas en Docker | Bajo |

---

## 12. Resumen de riesgos priorizados (estado actual)

| Prioridad | Vulnerabilidad | Impacto | Estado |
|---|---|---|---|
| 🔴 ~~Crítica~~ | ~~Sin rate limiting en auth~~ | - | ✅ **CORREGIDO** |
| 🔴 ~~Crítica~~ | ~~Sin límite de body HTTP~~ | - | ✅ **CORREGIDO** |
| 🟡 **Media** | `JWT_SECRET` placeholder en `.env.example` | Bajo | ⚠️ Persiste (documentación) |
| 🟠 ~~Alta~~ | ~~Sin `WriteTimeout` / `IdleTimeout`~~ | - | ✅ **CORREGIDO** |
| 🟠 ~~Alta~~ | ~~Contenedor Docker como root~~ | - | ✅ **CORREGIDO** |
| 🟠 ~~Alta~~ | ~~Headers de seguridad HTTP faltantes~~ | - | ✅ **CORREGIDO** |
| 🟠 **Alta** | `sslmode=disable` en DB | Alto | ⚠️ Persiste (solo local) |
| 🟡 ~~Media~~ | ~~JWT custom~~ | - | ✅ **CORREGIDO** (migrado) |
| 🟡 ~~Media~~ | ~~Sin refresh token rotation~~ | - | ✅ **CORREGIDO** |
| 🟡 ~~Media~~ | ~~Sin endpoint healthcheck~~ | - | ✅ **CORREGIDO** |
| 🟡 **Media** | Logger no estructurado / sin audit log | Bajo | ⚠️ Persiste |
| 🟢 **Baja** | Sin `nbf`, `jti` en JWT | Bajo | ⚠️ Persiste |
| 🟢 **Baja** | Sin MFA | Bajo | ⚠️ Persiste (contexto) |

### Estado actual

- **8 corregidas** de 13 vulnerabilidades identificadas inicialmente
- **5 persistentes**: 2 bajas (aceptables), 2 dependen de entorno/proceso, 1 de mejora menor

---

## 13. Fortalezas actuales (no tocar)

| Aspecto | Detalle |
|---|---|
| SQL parameterizado | Todas las queries usan placeholders. No hay SQL injection |
| Bcrypt con costo configurable | Contraseñas hasheadas con `golang.org/x/crypto` |
| CORS whitelist | Solo orígenes configurados, sin wildcards |
| Cache-Control: no-store | En todas las respuestas de handlers |
| Recoverer middleware | Previene caída del servidor por pánico |
| Context key tipo struct no exportado | Previene colisiones de contexto |
| Separación clean architecture | No hay lógica de negocio en handlers |
| Validación de input en use cases | Email, username, password mínima |
| JWT con verificación de algoritmo | Previene algorithm confusion attack |
| Contraseña no se devuelve en responses | Campo `PasswordHash` no se serializa |
| `.gitignore` incluye `.env` | No se trackean credenciales accidentalmente |

---

## 14. Próximos pasos sugeridos

### ✅ Ya implementados (no rehacer)

1. ~~Rate limiting~~ → ✅ `middleware/rate_limiter.go`
2. ~~Límite de body~~ → ✅ `middleware/max_body_bytes.go`
3. ~~Timeouts HTTP~~ → ✅ `cmd/api/main.go:97-100`
4. ~~Security headers~~ → ✅ `middleware/security_headers.go`
5. ~~Usuario no-root en Docker~~ → ✅ `Dockerfile:22-23`
6. ~~Healthcheck endpoint~~ → ✅ `routes/router.go:40-43`
7. ~~JWT con librería estándar~~ → ✅ `internal/infrastructure/security/jwt.go`
8. ~~Refresh token rotation~~ → ✅ `auth/auth.go:147-218`, `entities/refresh_token.go`, repositorio

### Pendientes

1. **Configurar `sslmode=require`** para conexiones de base de datos en producción (vía variable de entorno)
2. **Agregar structured logging** (`slog` o `zerolog`) con registro de eventos de seguridad (intentos de login fallidos, etc.)
3. **Validar en `config.go`** que `JWT_SECRET` no sea el placeholder `change-this-secret-in-production`
4. **Migrar a secrets management** (Docker secrets, HashiCorp Vault, o CI/CD env vars) para producción
5. **Agregar `nbf` y `jti`** en claims JWT si se requiere revocación individual de tokens

---

*Documento generado como parte del análisis de seguridad previo a la implementación de hardening.*
