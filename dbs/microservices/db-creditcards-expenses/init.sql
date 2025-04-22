CREATE TABLE credit_card_expenses (
    id SERIAL PRIMARY KEY,
    user_id INT,
    card_id INT,
    description VARCHAR(255) NOT NULL,
    amount DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
    purchase_date DATE NOT NULL,
    installment_count INT DEFAULT 1 CHECK (installment_count >= 1),
    category_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE
);
