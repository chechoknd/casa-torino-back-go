-- name: CreateIngredient :one
INSERT INTO ingredients (
    id,
    name,
    unit,
    average_cost,
    stock,
    minimum_stock,
    created_at,
    updated_at,
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetIngredientByID :one
SELECT *
FROM ingredients
WHERE id = $1;

-- name: ListIngredients :many
SELECT *
FROM ingredients
WHERE is_active = TRUE
ORDER BY created_at DESC;

-- name: UpdateIngredient :one
UPDATE ingredients
SET
    name = $2,
    unit = $3,
    average_cost = $4,
    stock = $5,
    minimum_stock = $6,
    updated_at = $7,
    is_active = $8
WHERE id = $1
RETURNING *;

-- name: DeactivateIngredient :execrows
UPDATE ingredients
SET
    is_active = FALSE,
    updated_at = $2
WHERE id = $1
  AND is_active = TRUE;
