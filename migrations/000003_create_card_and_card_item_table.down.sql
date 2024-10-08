-- Drop triggers
DROP TRIGGER IF EXISTS update_card_timestamp ON card;
DROP TRIGGER IF EXISTS update_card_item_timestamp ON card_item;

-- Drop foreign key constraints
ALTER TABLE card_item
DROP CONSTRAINT IF EXISTS fk_card_item_item_id;

ALTER TABLE card_item
DROP CONSTRAINT IF EXISTS fk_card_item_card_id;

ALTER TABLE card
DROP CONSTRAINT IF EXISTS fk_card_user_id;

-- Drop indexes
DROP INDEX IF EXISTS idx_card_id;
DROP INDEX IF EXISTS idx_card_user_id;
DROP INDEX IF EXISTS idx_card_item_item_id;
DROP INDEX IF EXISTS idx_card_item_card_id;

-- Drop tables
DROP TABLE IF EXISTS card_item;
DROP TABLE IF EXISTS card;
