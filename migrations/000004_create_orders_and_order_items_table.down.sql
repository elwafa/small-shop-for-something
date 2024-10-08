-- Drop triggers
DROP TRIGGER IF EXISTS update_orders_timestamp ON orders;
DROP TRIGGER IF EXISTS update_order_items_timestamp ON order_items;

-- Drop foreign key constraints
ALTER TABLE order_items
DROP CONSTRAINT IF EXISTS fk_order_items_item_id;

ALTER TABLE order_items
DROP CONSTRAINT IF EXISTS fk_order_items_order_id;

ALTER TABLE orders
DROP CONSTRAINT IF EXISTS fk_card_user_id;

-- Drop indexes
DROP INDEX IF EXISTS idx_order_id;
DROP INDEX IF EXISTS idx_order_created_by;
DROP INDEX IF EXISTS idx_order_items_item_id;
DROP INDEX IF EXISTS idx_order_items_order_id;

-- Drop tables
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS card;
