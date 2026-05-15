# Análisis Previo y Preguntas

> Fecha: 2026-05-13
> Propósito: Revisión de estado antes de nuevas tareas de desarrollo

---

## Resumen del estado actual

### Commits recientes en `feature/customer-panel-auth-roles`

| Commit | Mensaje |
|--------|---------|
| `52d1af5` | docs: add agent configurations, AGENTS.md, admin manual, and security analysis |
| `763138d` | docs: finalize customer panel hardening |
| `2747b5a` | feat: add public catalog and customer panel |
| `98747d3` | feat: add auth roles for customer panel |
| `029771f` | Merge branch 'feature/fix-vulnerabilities' into dev |

### Branches existentes sin merge a `dev`

- `bugfix/login-refresh-tokens-table`
- `feature/add-recipe-order`
- `feature/fix-vulnerabilities`
- `fix/recipes-payments-orders-display-consecutive`
- `fixBug/go-vet-recipe-unit`
- `fixBug/payment-order`
- `test/payment-module`

### Estado de sprints (IMPLEMENTATION_PLAN.md)

| Sprint | Estado |
|--------|--------|
| Sprint 1 - Backup, Auth y Roles | ✅ DONE |
| Sprint 2 - Customer Panel y Guest Mode | ✅ DONE |
| Sprint 3 - Integración, Migraciones y Hardening | ✅ DONE |

### Items bloqueados (DEVELOPMENT_STATUS.md)

| # | Item | Impacto |
|---|------|---------|
| 1 | Definir usuario o estrategia inicial para ADMIN | **Alto** - Sin admin inicial no se puede probar el flujo completo |
| 2 | Confirmar formato inicial de imagen de producto | **Medio** - Dependencia para frontend |
| 3 | Confirmar si refresh tokens quedan igual en esta fase | **Bajo** - Ya implementado con rotation |

---

## Observaciones y preguntas

### 1. Flujo de branching inconsistente

**AGENTS.md** dice:
- Bug-fixer: crear rama `fixBug/nombre-del-bug` desde `main`
- Feature: crear rama `feature/nombre-feature` desde `dev`

**Realidad en el repo:**
- Hay ramas `bugfix/login-refresh-tokens-table` (con formato distinto)
- La rama actual `feature/customer-panel-auth-roles` está basada en `dev`
- Hay ramas `fixBug/` y `fix/`

➡️ **Pregunta:** ¿El flujo correcto es desde `dev` siempre, o bugs van desde `main`? ¿El formato preferido es `fixBug/`, `bugfix/` o `fix/`?

### 2. Integración de branches a `dev`

Hay múltiples branches de fix/feature que parecen no haberse fusionado a `dev`:
- `bugfix/login-refresh-tokens-table`
- `fix/recipes-payments-orders-display-consecutive`
- `fixBug/go-vet-recipe-unit`
- `fixBug/payment-order`

➡️ **Pregunta:** ¿Estas ramas están pendientes de merge, ya no son relevantes, o se manejaron en otro lado?

### 3. `feature/customer-panel-auth-roles` sin mergear a `dev`

La rama actual está 3 commits adelante de `dev` pero no se ha fusionado. Tiene los sprints 1-3 completos.

➡️ **Pregunta:** ¿Debo mergear `feature/customer-panel-auth-roles` a `dev` antes de empezar nuevas tareas?

### 4. Usuario ADMIN inicial (BLOCKED #1)

El sistema soporta roles `ADMIN` y `CUSTOMER`, y el middleware requiere `ADMIN` para rutas administrativas. Pero no hay un usuario admin inicial creado ni una estrategia definida:
- ¿Migración de base de datos que cree el primer admin?
- ¿Script seed?
- ¿Asignación manual vía SQL?

➡️ **Pregunta:** ¿Cuál es la estrategia para crear/definir el usuario ADMIN inicial? ¿Credenciales esperadas?

### 5. Formato de imagen de producto (BLOCKED #2)

En la migración `000006_add_product_catalog_fields` se agregó `image_url` como texto, pero no hay confirmación del formato esperado:
- ¿URL externa (ej. S3, Cloudinary)?
- ¿Path interno?
- ¿Solo campo texto simple por ahora?

➡️ **Pregunta:** ¿El campo `image_url` actual como texto libre es suficiente, o se necesita cambiar antes de seguir?

### 6. Relación entre feature branch y dev

El plan de implementación y AGENTS.md usan `dev` como base. Sin embargo la rama actual está basada en `dev` y feature/customer-panel-auth-roles sería la rama de trabajo. Las nuevas features:
- ¿Se trabajan sobre `feature/customer-panel-auth-roles`?
- ¿Se mergea customer-panel-auth-roles a `dev` y se crean nuevas ramas desde `dev`?
- ¿Se crea una rama nueva desde `feature/customer-panel-auth-roles`?

### 7. Tests existentes

Los tests existentes pasan pero no está claro:
- ¿Hay tests de integración con base de datos?
- ¿Hay smoke tests automatizados ad-hoc?
- ¿Cobertura actual estimada?

### 8. Sprints completados

Los 3 sprints del plan están marcados como DONE. ¿Hay un siguiente plan/feature definido?:
- ¿Carrito de compras?
- ¿Integración WhatsApp?
- ¿Panel customer frontend?
- ¿Mejoras de seguridad/hardening?

---

## Archivos modificados en el commit actual

| Archivo | Propósito |
|---------|-----------|
| `.agents/AGENT_TESTER.md` | Config de agente QA/SDET |
| `.agents/tech-lead-review.md` | Config de agente tech lead reviewer |
| `AGENTS.md` | Reglas generales para agentes en este repo |
| `docs/admin-mvp-manual.md` | Manual/documentación para el MVP administrativo |
| `docs/vulnerabilities_risks.md` | Análisis de vulnerabilidades y seguridad |
