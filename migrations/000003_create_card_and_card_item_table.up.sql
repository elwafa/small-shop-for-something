-- create shopping card table
CREATE TABLE card (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- create card items table
CREATE TABLE card_item (
    id SERIAL PRIMARY KEY,
    card_id INT NOT NULL,
    item_id INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on the 'id' column for fast lookups
CREATE INDEX IF NOT EXISTS idx_card_id ON card (id);

-- Create an index on the 'user_id' column for fast lookups
CREATE INDEX IF NOT EXISTS idx_card_user_id ON card (user_id);

-- Create an index on the 'item_id' column for fast lookups
CREATE INDEX IF NOT EXISTS idx_card_item_item_id ON card_item (item_id);

-- Create an index on the 'card_id' column for fast lookups
CREATE INDEX IF NOT EXISTS idx_card_item_card_id ON card_item (card_id);

-- Create a foreign key constraint to the 'users' table on the 'user_id' column
ALTER TABLE card
ADD CONSTRAINT fk_card_user_id
FOREIGN KEY (user_id)
REFERENCES users (id)
ON DELETE CASCADE;

-- Create a foreign key constraint to the 'items' table on the 'item_id' column
ALTER TABLE card_item
ADD CONSTRAINT fk_card_item_item_id
FOREIGN KEY (item_id)
REFERENCES items (id)
ON DELETE CASCADE;

-- Create a foreign key constraint to the 'card' table on the 'card_id' column
ALTER TABLE card_item
ADD CONSTRAINT fk_card_item_card_id
FOREIGN KEY (card_id)
REFERENCES card (id)
ON DELETE CASCADE;

-- Create a trigger function to update the 'updated_at' column when a row is updated
CREATE OR REPLACE FUNCTION update_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
        NEW.updated_at = NOW();
RETURN NEW;
END;
    $$ LANGUAGE plpgsql;

CREATE TRIGGER update_card_timestamp
    BEFORE UPDATE ON card
    FOR EACH ROW
    EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_card_item_timestamp
    BEFORE UPDATE ON card_item
    FOR EACH ROW
    EXECUTE PROCEDURE update_timestamp();