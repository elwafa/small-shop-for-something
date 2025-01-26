-- Remove `category` and `colour` columns from items table
ALTER TABLE items
DROP COLUMN IF EXISTS category,
DROP COLUMN IF EXISTS colour;
