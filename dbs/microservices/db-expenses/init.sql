CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    user_id INT,
    description VARCHAR(255) NOT NULL,
    amount DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
    due_date DATE NOT NULL,
    paid BOOLEAN DEFAULT FALSE,
    category_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    user_id INT,
    expense_id INT,
    paid_at DATE NOT NULL,
    amount DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
    deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_expense_user ON expenses(user_id);
CREATE INDEX idx_expense_due_date ON expenses(user_id,due_date);
