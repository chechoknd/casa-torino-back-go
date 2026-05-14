# Development Status - Customer Panel

## TODO

- Preparar commit de Sprint 3.

## IN PROGRESS

- Preparar commit de Sprint 3.

## DONE

- Detectado `AGENTS.md` y flujo obligatorio del repo.
- Detectado auth existente con users, JWT y refresh tokens.
- Detectado middleware JWT existente.
- Detectadas rutas actuales protegidas globalmente.
- Detectado modelo product actual sin imagen ni visibilidad publica.
- Ejecutado `git status`.
- Cambiado a `dev`.
- Ejecutado `git pull origin dev`.
- Creada rama `feature/customer-panel-auth-roles`.
- Realizado backup de PostgreSQL en `backups/casa_torino_pre_customer_panel_20260513_174304.sql`.
- Documentada restauracion del backup en `docs/database-backup-restore.md`.
- Validada migracion actual: version `4`, dirty `false`.
- Ejecutado `go fmt ./...`.
- Ejecutado `go vet ./...`.
- Ejecutado `go test ./...`.
- Revisado auth/users/JWT existentes.
- Agregada migracion `000005_add_user_roles`.
- Aplicada migracion local: version `5`, dirty `false`.
- Usuarios existentes migrados a rol `CUSTOMER`.
- Agregado rol a entidad `User`.
- Incluido rol en DTOs de auth.
- Incluido rol en JWT.
- Expuesto rol en middleware auth.
- Creado middleware de roles.
- Protegidas rutas existentes con rol `ADMIN`.
- Ajustados tests de auth/JWT/middleware/routes.
- Validada reversibilidad de migracion `000005` con `migrate-down` y `migrate-up`.
- Ejecutado `go fmt ./...` despues de cambios de Sprint 1.
- Ejecutado `go vet ./...` despues de cambios de Sprint 1.
- Ejecutado `go test ./...` despues de cambios de Sprint 1.
- Backend reconstruido con Docker para smoke tests de Sprint 1.
- Smoke test `GET /health`: `200 OK`.
- Smoke test register/login: respuesta incluye `role`.
- Smoke test ruta admin sin token: `401 UNAUTHORIZED`.
- Smoke test ruta admin con token `CUSTOMER`: `403 FORBIDDEN`.
- Smoke test ruta admin con token `ADMIN`: `200 OK`.
- Usuario temporal de smoke test eliminado de la base local.
- Sprint 1 - Backup, Auth y Roles completado.
- Revisado products/rutas actuales para catalogo publico.
- Agregada migracion `000006_add_product_catalog_fields`.
- Aplicada migracion local: version `6`, dirty `false`.
- Backfill de productos: activos como publicos, inactivos como no publicos.
- Agregados campos `image_url` e `is_public` a products.
- Actualizados entidad, DTOs, repository, mapper y sqlc de products.
- Agregados endpoints publicos `/public/products`, `/public/products/{id}` y `/public/product-categories`.
- Ajustados tests de products/catalogo.
- Agregados endpoints customer `/customer/profile` y `/customer/orders`.
- Panel cliente usa email del usuario autenticado para resolver customer existente.
- Validado guest catalog sin JWT.
- Validado guest sin acceso a `/orders` ni `/customer/profile`.
- Validado customer con acceso a `/customer/profile` y `/customer/orders`.
- Validada reversibilidad de migracion `000006` con `migrate-down` y `migrate-up`.
- Ejecutado `go fmt ./...` despues de cambios de Sprint 2.
- Ejecutado `go vet ./...` despues de cambios de Sprint 2.
- Ejecutado `go test ./...` despues de cambios de Sprint 2.
- Sprint 2 - Customer Panel y Guest Mode completado.
- Backup final generado en `backups/casa_torino_post_customer_panel_20260513_190909.sql`.
- Validada migracion final: version `6`, dirty `false`.
- Actualizados README y docs con roles, rutas publicas, rutas customer y backup final.
- Validado flujo guest: catalogo publico `200 OK`, ruta privada `401 UNAUTHORIZED`.
- Validado flujo admin: login con rol `ADMIN`, products/orders/payments `200 OK`.
- Validado flujo customer: profile/orders `200 OK`, ruta admin `403 FORBIDDEN`.
- Usuarios temporales de smoke test final eliminados.
- Ejecutado `go fmt ./...` despues de cambios de Sprint 3.
- Ejecutado `go vet ./...` despues de cambios de Sprint 3.
- Ejecutado `go test ./...` despues de cambios de Sprint 3.
- Sprint 3 - Integracion, Migraciones y Hardening completado.

## BLOCKED

- Definir usuario o estrategia inicial para `ADMIN`.
- Confirmar formato inicial de imagen de producto: URL externa, path interno o storage futuro.
- Confirmar si refresh tokens quedan igual en esta fase.
