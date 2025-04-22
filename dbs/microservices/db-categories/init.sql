CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    user_id INT,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(9) NOT NULL DEFAULT '#00000000' CHECK (color ~ '^#[0-9A-Fa-f]{8}$'),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    active BOOLEAN DEFAULT TRUE,
    UNIQUE (user_id, name)
);
