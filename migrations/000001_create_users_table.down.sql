-- Drop the trigger
DROP TRIGGER IF EXISTS update_users_timestamp ON users;

-- Drop the function
DROP FUNCTION IF EXISTS update_timestamp;

-- Drop the table
DROP TABLE IF EXISTS users;