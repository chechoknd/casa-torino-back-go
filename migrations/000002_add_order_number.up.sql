CREATE SEQUENCE IF NOT EXISTS orders_order_number_seq;

ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS order_number BIGINT;

ALTER TABLE orders
    ALTER COLUMN order_number SET DEFAULT nextval('orders_order_number_seq');

WITH ordered_rows AS (
    SELECT id, nextval('orders_order_number_seq') AS generated_number
    FROM orders
    WHERE order_number IS NULL
    ORDER BY created_at, id
)
UPDATE orders AS o
SET order_number = ordered_rows.generated_number
FROM ordered_rows
WHERE o.id = ordered_rows.id;

SELECT setval(
    'orders_order_number_seq',
    COALESCE((SELECT MAX(order_number) FROM orders), 1),
    false
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_order_number_unique
    ON orders (order_number);
