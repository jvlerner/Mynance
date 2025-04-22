CREATE TABLE credit_cards (
    id SERIAL PRIMARY KEY,
    user_id INT,
    name VARCHAR(100) NOT NULL,
    bank VARCHAR(50) NOT NULL,
    limit_amount DECIMAL(10,2) NOT NULL CHECK (limit_amount >= 0),
    due_day INT NOT NULL CHECK (due_day BETWEEN 1 AND 31),
    active BOOLEAN DEFAULT TRUE
);

-- Índices para performance
CREATE INDEX idx_card_user ON credit_cards(user_id);
