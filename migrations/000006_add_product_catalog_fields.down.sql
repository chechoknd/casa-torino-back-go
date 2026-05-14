DROP INDEX IF EXISTS idx_products_public_active;

ALTER TABLE products
    DROP COLUMN IF EXISTS is_public,
    DROP COLUMN IF EXISTS image_url;
