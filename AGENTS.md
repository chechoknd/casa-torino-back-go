# AGENTS.md

## Rol general

Este es un backend en Go. Todo agente debe actuar como desarrollador senior cuidadoso.

## Reglas obligatorias

- No hacer cambios grandes sin explicar el plan.
- No modificar arquitectura sin justificarlo.
- No romper APIs existentes.
- No ignorar errores.
- No hacer cambios cosméticos innecesarios.
- Antes de tocar código, revisar estructura existente.
- Mantener estilo idiomático de Go.

## Comandos de verificación

Ejecutar antes de dar por terminado:

go fmt ./...
go test ./...
go vet ./...

## Arquitectura esperada

- Mantener separación entre handlers, services/usecases, repositories e infraestructura.
- No meter lógica de negocio directamente en handlers.
- No acceder a base de datos desde capas superiores si ya existe repository.
- Preferir interfaces pequeñas.
- Manejar errores explícitamente.
- Usar context.Context cuando aplique.

## Criterios de finalización

Una tarea solo está terminada si:

- Compila.
- Tests pasan.
- No hay errores de go vet.
- El cambio es pequeño y revisable.
- Se explica qué archivos cambiaron y por qué.

## Modo bug-fixer

Cuando el usuario diga "usa el modo bug-fixer", actúa como desarrollador senior especializado en corrección de bugs.

### Reglas obligatorias

- Antes de modificar código, revisar el contexto del bug.
- Crear siempre una rama nueva desde `main` actualizada.
- La rama debe usar el formato:

fixBug/nombre-del-bug

- No mezclar refactors grandes con correcciones pequeñas.
- No cambiar arquitectura sin necesidad.
- Respetar la arquitectura limpia existente del proyecto.
- Tomar en cuenta todas las reglas generales definidas en este `AGENTS.md`.
- Hacer el cambio mínimo necesario para resolver el bug.
- Agregar o ajustar tests solo si aplica.
- Ejecutar al final:

go fmt ./...
go vet ./...
go test ./...

### Flujo obligatorio

1. Verificar estado actual del repo
git status

2. Cambiar a dev
git checkout dev

3. Actualizar dev
git pull origin dev

4. Crear rama desde dev
git checkout -b bugfix/nombre-del-bug

Y para features:

git checkout dev
git pull origin dev

git checkout -b feature/nombre-feature

5. Reproducir o confirmar el bug.

6. Aplicar la corrección mínima.

7. Ejecutar validaciones:

go fmt ./...
go vet ./...
go test ./...

8. Mostrar resumen final con:

- Rama creada.
- Archivos modificados.
- Causa del bug.
- Solución aplicada.
- Resultado de comandos.
- Riesgos restantes.
