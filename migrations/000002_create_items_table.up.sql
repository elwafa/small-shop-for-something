CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2),
    picture TEXT,
    status VARCHAR(10),
    receive VARCHAR(100),
    user_id INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
        NEW.updated_at = NOW();
RETURN NEW;
END;
    $$ LANGUAGE plpgsql;

CREATE TRIGGER update_items_timestamp
    BEFORE UPDATE ON items
    FOR EACH ROW
    EXECUTE PROCEDURE update_timestamp();


-- Create an index on the 'id' column for fast lookups
CREATE INDEX IF NOT EXISTS idx_items_id ON items (id);

-- Create an index on the 'name' column for fast lookups
CREATE INDEX IF NOT EXISTS idx_items_name ON items (name);

-- Create an index on the 'status' column for fast lookups
CREATE INDEX IF NOT EXISTS idx_items_status ON items (status);

-- Create an index on the 'receive' column for fast lookups
CREATE INDEX IF NOT EXISTS idx_items_receive ON items (receive);

-- Create an index on the 'user_id' column for fast lookups
CREATE INDEX IF NOT EXISTS idx_items_user_id ON items (user_id);

-- Create a foreign key constraint to the 'users' table on the 'user_id' column

ALTER TABLE items
ADD CONSTRAINT fk_items_user_id
FOREIGN KEY (user_id)
REFERENCES users (id)
ON DELETE CASCADE;