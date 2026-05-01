-- name: CreateCustomer :one
INSERT INTO customers (
    id,
    full_name,
    phone,
    email,
    customer_type,
    created_at,
    updated_at,
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetCustomerByID :one
SELECT *
FROM customers
WHERE id = $1;

-- name: GetCustomerByEmail :one
SELECT *
FROM customers
WHERE LOWER(email) = LOWER($1);

-- name: ListCustomers :many
SELECT *
FROM customers
WHERE is_active = TRUE
ORDER BY created_at DESC;

-- name: UpdateCustomer :one
UPDATE customers
SET
    full_name = $2,
    phone = $3,
    email = $4,
    customer_type = $5,
    updated_at = $6,
    is_active = $7
WHERE id = $1
RETURNING *;

-- name: DeactivateCustomer :execrows
UPDATE customers
SET
    is_active = FALSE,
    updated_at = $2
WHERE id = $1
  AND is_active = TRUE;
