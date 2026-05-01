-include .env
export

MIGRATIONS_DATABASE_URL ?= $(DATABASE_URL)

.PHONY: run down migrate-up migrate-down sqlc test db-counts db-shell

run:
	docker compose up --build

down:
	docker compose down

migrate-up:
	go run -tags 'postgres file' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3 -path migrations -database "$(MIGRATIONS_DATABASE_URL)" up

migrate-down:
	go run -tags 'postgres file' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3 -path migrations -database "$(MIGRATIONS_DATABASE_URL)" down 1

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0 generate

test:
	go test ./...

db-counts:
	docker compose exec -T db psql -U user -d casa_torino -c "SELECT 'customers' AS table_name, COUNT(*) AS total FROM customers UNION ALL SELECT 'products', COUNT(*) FROM products UNION ALL SELECT 'ingredients', COUNT(*) FROM ingredients UNION ALL SELECT 'recipes', COUNT(*) FROM recipes UNION ALL SELECT 'recipe_items', COUNT(*) FROM recipe_items UNION ALL SELECT 'orders', COUNT(*) FROM orders UNION ALL SELECT 'order_items', COUNT(*) FROM order_items UNION ALL SELECT 'payments', COUNT(*) FROM payments ORDER BY table_name;"

db-shell:
	docker compose exec db psql -U user -d casa_torino
