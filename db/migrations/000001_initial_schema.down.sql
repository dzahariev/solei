BEGIN;

-- Drop the triggers
DROP TRIGGER set_created_at_on_orderitems ON orderitems;
DROP TRIGGER set_updated_at_on_orderitems ON orderitems;
DROP TRIGGER set_created_at_on_orders ON orders;
DROP TRIGGER set_updated_at_on_orders ON orders;
DROP TRIGGER set_created_at_on_meals ON meals;
DROP TRIGGER set_updated_at_on_meals ON meals;
DROP TRIGGER set_created_at_on_categories ON categories;
DROP TRIGGER set_updated_at_on_categories ON categories;
DROP TRIGGER set_created_at_on_addresses ON addresses;
DROP TRIGGER set_updated_at_on_addresses ON addresses;
DROP TRIGGER set_created_at_on_users ON users;
DROP TRIGGER set_updated_at_on_users ON users;

-- Drop the helper funcitons
DROP FUNCTION set_created_at();
DROP FUNCTION set_updated_at();

-- Drop tables
DROP TABLE orderitems;
DROP TABLE orders;
DROP TABLE meals;
DROP TABLE categories;
DROP TABLE addresses;
DROP TABLE users;

COMMIT;