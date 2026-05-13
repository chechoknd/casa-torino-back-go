# Casa Torino - Customer Panel Implementation Plan

## Resumen corto

Extender el MVP actual para soportar roles, login con respuesta de rol, modo invitado, catalogo publico para clientes y proteccion explicita de rutas admin. El cambio debe ser incremental sobre la arquitectura hexagonal existente en Go, PostgreSQL y Docker, reutilizando los modulos actuales de `users`, `auth`, JWT, middleware, products, orders y customers.

No se replantea la arquitectura. No se redisenan modulos existentes. La prioridad es mantener el MVP actual funcionando y agregar permisos y rutas publicas de forma pequena y verificable.

## Reglas del repo que aplican

- Leer y respetar `AGENTS.md` antes de implementar.
- Para feature: trabajar desde `dev`, actualizar con `git pull origin dev` y crear rama `feature/customer-panel-auth-roles`.
- Mantener separacion entre handlers, usecases, repositories e infraestructura.
- No meter logica de negocio directamente en handlers.
- No romper APIs existentes.
- Antes de cerrar cada sprint ejecutar:
  - `go fmt ./...`
  - `go test ./...`
  - `go vet ./...`

## Modulos impactados

- `users`: agregar rol y relacion opcional con customer si aplica.
- `customers`: vincular cliente autenticado con customer existente o nuevo.
- `auth`: login/register deben devolver rol; validar estado activo.
- `JWT`: incluir rol en claims para autorizacion.
- `middleware`: agregar middleware de roles sobre el JWT actual.
- `routes`: separar rutas publicas, rutas customer y rutas admin.
- `products`: agregar imagenes, visibilidad publica y datos de catalogo.
- `orders`: impedir compras en modo guest; permitir compras solo a customer/admin segun flujo actual.
- `payments`: mantener protegido; no exponer a guest.
- `migrations`: agregar cambios compatibles y reversibles.
- `docs`: documentar backup y restauracion.

## Orden recomendado de implementacion

1. Backup y validacion de seguridad de datos.
2. Revisar estado real de auth/users/JWT existente.
3. Agregar roles mínimos (`ADMIN`, `CUSTOMER`) en usuarios.
4. Ajustar JWT y respuesta de login para incluir rol.
5. Agregar middleware de roles.
6. Separar rutas publicas de rutas protegidas.
7. Agregar campos publicos de productos para catalogo.
8. Agregar endpoints publicos de catalogo.
9. Proteger rutas admin explicitamente.
10. Validar guest mode desde el contrato de rutas: guest no necesita JWT y solo accede a rutas publicas.

## Riesgos importantes

- Ya existen `users`, `auth`, JWT y refresh tokens; crear duplicados romperia el MVP.
- Cambiar todas las rutas protegidas de golpe puede romper clientes actuales.
- Relacionar users y customers puede afectar datos existentes si se fuerza `NOT NULL` sin migracion gradual.
- Agregar roles sin backfill puede dejar usuarios existentes sin permisos.
- Imagenes de productos deben iniciar como URL o metadata simple; no conviene agregar storage complejo en esta fase.
- Guest no debe ser persistido como usuario real salvo que haya una necesidad clara.
- Endpoints publicos deben filtrar solo productos activos y visibles.

## Dependencias

- Acceso a base de datos para backup y prueba de migraciones.
- Definir credenciales/usuario inicial `ADMIN` para el ambiente actual.
- Confirmar si customers actuales deben mapearse automaticamente a users por email.
- Confirmar URL/base path esperado por el frontend para catalogo publico.
- Confirmar si imagen de producto sera URL externa, path interno o campo simple inicial.

## Sprint 1 - Backup, Auth y Roles

**Estado:** DONE

### Features

- Backup obligatorio antes de migraciones.
- Roles `ADMIN` y `CUSTOMER`.
- Login devuelve rol.
- JWT incluye rol.
- Middleware de roles.
- Rutas admin protegidas explicitamente.

### Tasks

- [x] Revisar `AGENTS.md` y confirmar flujo de feature branch desde `dev`.
- [x] Ejecutar `git status`.
- [x] Cambiar a `dev`.
- [x] Ejecutar `git pull origin dev`.
- [x] Crear rama `feature/customer-panel-auth-roles`.
- [x] Realizar backup de PostgreSQL antes de tocar migraciones.
- [x] Documentar comando de backup usado y ubicacion del archivo.
- [x] Documentar restauracion con `pg_restore` o `psql`, segun formato del backup.
- [x] Probar migraciones actuales en base local o entorno de prueba.
- [x] Revisar implementacion existente de `users`, `auth`, JWT y refresh tokens.
- [x] Crear migracion para agregar `role` a `users` con default seguro.
- [x] Agregar backfill para usuarios existentes.
- [x] Actualizar entidad `User` con rol.
- [x] Actualizar repository postgres/sqlc si aplica.
- [x] Actualizar DTOs de auth para incluir rol en login/register/refresh.
- [x] Actualizar JWT service para emitir y verificar rol.
- [x] Actualizar middleware auth para exponer rol en contexto.
- [x] Crear middleware `RequireRole`.
- [x] Proteger rutas admin actuales: customers, products admin, ingredients, recipes, orders admin y payments.
- [x] Mantener compatibilidad con endpoints existentes donde sea posible.
- [x] Agregar o ajustar tests de auth, JWT, middleware y router.
- [x] Ejecutar `go fmt ./...`.
- [x] Ejecutar `go vet ./...`.
- [x] Ejecutar `go test ./...`.

### Nota operativa

- La estrategia para crear o asignar el primer `ADMIN` queda pendiente de decision. Mientras tanto, el sistema soporta el rol y puede asignarse por migracion/manual SQL controlado si el entorno lo requiere.

### Resultado esperado

- Usuarios tienen rol.
- Login responde rol.
- JWT transporta rol.
- Rutas admin requieren `ADMIN`.
- Usuarios sin rol migran sin romper datos existentes.
- Hay backup y documento de restauracion.

## Sprint 2 - Customer Panel y Guest Mode

### Features

- Rutas publicas para catalogo.
- Guest mode por acceso anonimo a rutas publicas.
- Productos con imagen y visibilidad publica.
- Base para panel cliente.

### Tasks

- Crear migracion para productos:
  - `image_url`
  - `is_public`
  - campos simples de catalogo solo si hacen falta.
- Agregar backfill: productos activos pueden iniciar como publicos si el negocio lo acepta.
- Actualizar entidad, DTOs, repository y mapper de products.
- Agregar endpoint publico para listar productos visibles.
- Agregar endpoint publico para detalle de producto visible.
- Agregar endpoint publico de categorias usando `product_type` existente, salvo que ya exista otra fuente.
- Separar rutas publicas de `/products` admin si es necesario, por ejemplo `/public/products`.
- Validar que guest pueda ver catalogo sin JWT.
- Validar que guest no pueda crear orders, ver historial ni acceder a admin.
- Agregar endpoints customer protegidos por `CUSTOMER` para perfil/historial solo si el MVP actual ya tiene datos suficientes.
- Relacionar `users` con `customers` de forma gradual si el panel necesita historial por customer.
- Agregar tests de rutas publicas, permisos guest y permisos customer.
- Ejecutar `go fmt ./...`.
- Ejecutar `go vet ./...`.
- Ejecutar `go test ./...`.

### Resultado esperado

- El catalogo puede consumirse sin login.
- Guest queda representado por ausencia de JWT, no por un usuario persistido.
- Customer puede usar rutas cliente protegidas.
- Admin conserva acceso a operaciones internas.

## Sprint 3 - Integracion, Migraciones y Hardening

### Features

- Validacion completa del flujo MVP.
- Migraciones seguras.
- Checklist de release corto.

### Tasks

- Ejecutar backup nuevo antes de probar migraciones finales.
- Probar aplicar y revertir migraciones en entorno local.
- Validar que customers, products, ingredients, recipes, orders y payments existentes siguen funcionando.
- Probar flujo admin: login, crear producto, actualizar producto, listar orders, payments.
- Probar flujo customer: register/login, ver catalogo, crear order si aplica al MVP.
- Probar flujo guest: ver catalogo y recibir unauthorized/forbidden en rutas protegidas.
- Revisar respuestas de error para distinguir `401 Unauthorized` y `403 Forbidden`.
- Actualizar README o docs minimos con roles, backup y rutas nuevas.
- Ejecutar `go fmt ./...`.
- Ejecutar `go vet ./...`.
- Ejecutar `go test ./...`.

### Resultado esperado

- MVP actual sigue estable.
- Nueva fase queda lista para conectar frontend cliente.
- Riesgo de perdida de datos reducido por backup probado y migraciones reversibles.
