-- name: CreateProduct :one
INSERT INTO products (
    id,
    name,
    description,
    product_type,
    base_price,
    cost_price,
    image_url,
    is_public,
    is_active,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: GetProductByID :one
SELECT *
FROM products
WHERE id = $1;

-- name: GetPublicProductByID :one
SELECT *
FROM products
WHERE id = $1
  AND is_active = TRUE
  AND is_public = TRUE;

-- name: ListProducts :many
SELECT *
FROM products
WHERE is_active = TRUE
ORDER BY created_at DESC;

-- name: ListPublicProducts :many
SELECT *
FROM products
WHERE is_active = TRUE
  AND is_public = TRUE
ORDER BY created_at DESC;

-- name: ListPublicProductTypes :many
SELECT DISTINCT product_type
FROM products
WHERE is_active = TRUE
  AND is_public = TRUE
ORDER BY product_type ASC;

-- name: UpdateProduct :one
UPDATE products
SET
    name = $2,
    description = $3,
    product_type = $4,
    base_price = $5,
    cost_price = $6,
    image_url = $7,
    is_public = $8,
    is_active = $9,
    updated_at = $10
WHERE id = $1
RETURNING *;

-- name: DeactivateProduct :execrows
UPDATE products
SET
    is_active = FALSE,
    updated_at = $2
WHERE id = $1
  AND is_active = TRUE;
