-- name: CreateRecipe :one
INSERT INTO recipes (
    id,
    product_id,
    name,
    portions,
    created_at,
    updated_at,
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: CreateRecipeItem :one
INSERT INTO recipe_items (
    id,
    recipe_id,
    ingredient_id,
    quantity,
    unit
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetRecipeByID :one
SELECT *
FROM recipes
WHERE id = $1;

-- name: GetRecipeByProductID :one
SELECT *
FROM recipes
WHERE product_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: GetRecipeItemsByRecipeID :many
SELECT *
FROM recipe_items
WHERE recipe_id = $1
ORDER BY id;

-- name: ListRecipes :many
SELECT *
FROM recipes
WHERE is_active = TRUE
ORDER BY created_at DESC;

-- name: UpdateRecipe :one
UPDATE recipes
SET
    product_id = $2,
    name = $3,
    portions = $4,
    updated_at = $5,
    is_active = $6
WHERE id = $1
RETURNING *;

-- name: UpdateRecipeItem :one
UPDATE recipe_items
SET
    ingredient_id = $2,
    quantity = $3,
    unit = $4
WHERE id = $1
RETURNING *;

-- name: DeactivateRecipe :exec
UPDATE recipes
SET
    is_active = FALSE,
    updated_at = $2
WHERE id = $1;
