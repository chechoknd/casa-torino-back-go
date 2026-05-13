# Database Backup and Restore

## Backup inicial Customer Panel

Backup generado antes de iniciar cambios de desarrollo:

```bash
backups/casa_torino_pre_customer_panel_20260513_174304.sql
```

Estado de la base al momento del backup:

- `schema_migrations.version`: `4`
- `schema_migrations.dirty`: `false`
- `customers`: `13`
- `products`: `12`
- `ingredients`: `17`
- `recipes`: `12`
- `orders`: `8`
- `payments`: `6`

Comando usado:

```bash
docker compose exec -T db pg_dump -U user -d casa_torino --clean --if-exists --no-owner --no-privileges > backups/casa_torino_pre_customer_panel_20260513_174304.sql
```

## Restaurar backup local

Con los contenedores activos:

```bash
docker compose exec -T db psql -U user -d casa_torino < backups/casa_torino_pre_customer_panel_20260513_174304.sql
```

Luego validar migraciones:

```bash
docker compose exec -T db psql -U user -d casa_torino -c "SELECT version, dirty FROM schema_migrations;"
```

Y validar conteos basicos:

```bash
make db-counts
```
