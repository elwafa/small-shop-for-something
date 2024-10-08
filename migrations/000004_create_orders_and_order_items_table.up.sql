-- create shopping card table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    created_by INT NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- create card items table
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    item_id INT NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on the 'id' column for fast lookups
CREATE INDEX IF NOT EXISTS idx_orders_id ON orders (id);

-- Create an index on the 'user_id' column for fast lookups
CREATE INDEX IF NOT EXISTS idx_order_created_by ON orders (created_by);

-- Create an index on the 'item_id' column for fast lookups
CREATE INDEX IF NOT EXISTS idx_order_items_item_id ON order_items (item_id);

-- Create an index on the 'card_id' column for fast lookups
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items (order_id);


-- Create a foreign key constraint to the 'users' table on the 'user_id' column
ALTER TABLE orders
ADD CONSTRAINT fk_orders_created_by
FOREIGN KEY (created_by)
REFERENCES users (id)
ON DELETE CASCADE;

-- Create a foreign key constraint to the 'items' table on the 'item_id' column
ALTER TABLE order_items
ADD CONSTRAINT fk_order_items_item_id
FOREIGN KEY (item_id)
REFERENCES items (id)
ON DELETE CASCADE;

-- Create a foreign key constraint to the 'card' table on the 'card_id' column
ALTER TABLE order_items
ADD CONSTRAINT fk_order_items_order_id
FOREIGN KEY (order_id)
REFERENCES orders (id)
ON DELETE CASCADE;

-- Create a trigger function to update the 'updated_at' column when a row is updated
CREATE OR REPLACE FUNCTION update_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
        NEW.updated_at = NOW();
RETURN NEW;
END;
    $$ LANGUAGE plpgsql;

CREATE TRIGGER update_orders_timestamp
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_order_items_timestamp
    BEFORE UPDATE ON order_items
    FOR EACH ROW
    EXECUTE PROCEDURE update_timestamp();