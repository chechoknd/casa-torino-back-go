-- name: CreateOrder :one
INSERT INTO orders (
    id,
    customer_id,
    order_number,
    status,
    subtotal,
    discount,
    total,
    created_at,
    updated_at
) VALUES (
    $1, $2, DEFAULT, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (
    id,
    order_id,
    product_id,
    quantity,
    unit_price,
    subtotal
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetOrderByID :one
SELECT *
FROM orders
WHERE id = $1;

-- name: GetOrderItemsByOrderID :many
SELECT *
FROM order_items
WHERE order_id = $1
ORDER BY id;

-- name: ListOrdersByCustomerID :many
SELECT *
FROM orders
WHERE customer_id = $1
ORDER BY created_at DESC;

-- name: ListOrders :many
SELECT *
FROM orders
ORDER BY created_at DESC;

-- name: UpdateOrder :one
UPDATE orders
SET
    customer_id = $2,
    status = $3,
    subtotal = $4,
    discount = $5,
    total = $6,
    updated_at = $7
WHERE id = $1
RETURNING *;

-- name: UpdateOrderItem :one
UPDATE order_items
SET
    product_id = $2,
    quantity = $3,
    unit_price = $4,
    subtotal = $5
WHERE id = $1
RETURNING *;

-- name: UpdateOrderStatus :one
UPDATE orders
SET
    status = $2,
    updated_at = $3
WHERE id = $1
RETURNING *;
