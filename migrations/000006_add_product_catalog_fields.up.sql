ALTER TABLE products
    ADD COLUMN IF NOT EXISTS image_url TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS is_public BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE products
SET is_public = TRUE
WHERE is_active = TRUE;

CREATE INDEX IF NOT EXISTS idx_products_public_active
    ON products (product_type, created_at DESC)
    WHERE is_active = TRUE AND is_public = TRUE;
