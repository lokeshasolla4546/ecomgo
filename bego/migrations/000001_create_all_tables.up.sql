CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    role VARCHAR NOT NULL
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    price INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    image TEXT
);

CREATE TABLE carts (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    items JSON NOT NULL,
    total INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
