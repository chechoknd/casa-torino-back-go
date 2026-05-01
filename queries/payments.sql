-- name: CreatePayment :one
INSERT INTO payments (
    id,
    order_id,
    amount,
    method,
    status,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetPaymentByID :one
SELECT *
FROM payments
WHERE id = $1;

-- name: GetPaymentsByOrderID :many
SELECT *
FROM payments
WHERE order_id = $1
ORDER BY created_at DESC;

-- name: ListPayments :many
SELECT *
FROM payments
ORDER BY created_at DESC;

-- name: UpdatePayment :one
UPDATE payments
SET
    order_id = $2,
    amount = $3,
    method = $4,
    status = $5,
    updated_at = $6
WHERE id = $1
RETURNING *;
