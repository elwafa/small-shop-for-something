-- Add `category` and `colour` columns to items table
ALTER TABLE items
ADD COLUMN category VARCHAR(255),
ADD COLUMN colour VARCHAR(255);
