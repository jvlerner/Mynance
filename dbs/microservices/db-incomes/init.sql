CREATE TABLE incomes (
    id SERIAL PRIMARY KEY,
    user_id INT ,
    description VARCHAR(255) NOT NULL,
    amount DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
    received_at DATE NOT NULL,
    is_recurring BOOLEAN DEFAULT FALSE,
    deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_expense_user ON expenses(user_id);
CREATE INDEX idx_expense_due_date ON expenses(user_id,due_date);
