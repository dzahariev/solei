BEGIN;

-- Function that sets created_at field
CREATE OR REPLACE FUNCTION set_created_at()
    RETURNS TRIGGER
AS
$$
BEGIN
    NEW.created_at = NOW();
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

-- Function that sets updated_at field
CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER
AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

-- Table for users
CREATE TABLE users(
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    prefered_user_name VARCHAR(1024) NOT NULL,
    given_name VARCHAR(1024) NOT NULL,
    family_name VARCHAR(1024) NOT NULL,
    email VARCHAR(1024) NOT NULL
);

-- Trigger that sets created_at on users
CREATE TRIGGER set_created_at_on_users
    BEFORE INSERT 
    ON users
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

-- Trigger that sets created_at on users
CREATE TRIGGER set_updated_at_on_users
    BEFORE INSERT OR UPDATE 
    ON users
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Table for addresses
CREATE TABLE addresses(
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    contry VARCHAR(1024) NOT NULL,
    city VARCHAR(1024) NOT NULL,
    street VARCHAR(1024) NOT NULL,
    phone VARCHAR(1024) NOT NULL,
    user_id uuid NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Trigger that sets created_at on addresses
CREATE TRIGGER set_created_at_on_addresses
    BEFORE INSERT 
    ON addresses
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

-- Trigger that sets created_at on addresses
CREATE TRIGGER set_updated_at_on_addresses
    BEFORE INSERT OR UPDATE 
    ON addresses
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Table for categories
CREATE TABLE categories(
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    name VARCHAR(1024) NOT NULL
);

-- Trigger that sets created_at on categories
CREATE TRIGGER set_created_at_on_categories
    BEFORE INSERT 
    ON categories
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

-- Trigger that sets created_at on categories
CREATE TRIGGER set_updated_at_on_categories
    BEFORE INSERT OR UPDATE 
    ON categories
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Table for meals
CREATE TABLE meals(
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    name VARCHAR(1024) NOT NULL,
    description VARCHAR(1024) NOT NULL,
    cost real NOT NULL,
    category_id uuid NOT NULL,
    CONSTRAINT fk_category FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- Trigger that sets created_at on meals
CREATE TRIGGER set_created_at_on_meals
    BEFORE INSERT 
    ON meals
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

-- Trigger that sets created_at on meals
CREATE TRIGGER set_updated_at_on_meals
    BEFORE INSERT OR UPDATE 
    ON meals
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Table for orders
CREATE TABLE orders(
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    price real NOT NULL,
    status VARCHAR(1024) NOT NULL,
    user_id uuid NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Trigger that sets created_at on orders
CREATE TRIGGER set_created_at_on_orders
    BEFORE INSERT 
    ON orders
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

-- Trigger that sets created_at on orders
CREATE TRIGGER set_updated_at_on_orders
    BEFORE INSERT OR UPDATE 
    ON orders
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Table for orderitems
CREATE TABLE orderitems(
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    amount integer,
    comment VARCHAR(1024),
    meal_id uuid NOT NULL,
    order_id uuid NOT NULL,
    user_id uuid NOT NULL,
    CONSTRAINT fk_meal FOREIGN KEY(meal_id) REFERENCES meals(id) ON DELETE CASCADE,
    CONSTRAINT fk_order FOREIGN KEY(order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Trigger that sets created_at on orderitems
CREATE TRIGGER set_created_at_on_orderitems
    BEFORE INSERT 
    ON orderitems
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

-- Trigger that sets created_at on orderitems
CREATE TRIGGER set_updated_at_on_orderitems
    BEFORE INSERT OR UPDATE 
    ON orderitems
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

COMMIT;