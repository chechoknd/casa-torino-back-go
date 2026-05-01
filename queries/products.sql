-- name: CreateProduct :one
INSERT INTO products (
    id,
    name,
    description,
    product_type,
    base_price,
    cost_price,
    is_active,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetProductByID :one
SELECT *
FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT *
FROM products
WHERE is_active = TRUE
ORDER BY created_at DESC;

-- name: UpdateProduct :one
UPDATE products
SET
    name = $2,
    description = $3,
    product_type = $4,
    base_price = $5,
    cost_price = $6,
    is_active = $7,
    updated_at = $8
WHERE id = $1
RETURNING *;

-- name: DeactivateProduct :execrows
UPDATE products
SET
    is_active = FALSE,
    updated_at = $2
WHERE id = $1
  AND is_active = TRUE;
