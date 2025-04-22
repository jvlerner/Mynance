package db

import (
	"database/sql"

	"github.com/jvlerner/my-finance-api/pkg/postgres"
)

// GetExpensesByUser retrieves all expenses for a specific user
func GetExpensesByUser(userID int) ([]postgres.Expense, error) {
	rows, err := postgres.DB.Query("SELECT id, user_id, description, amount, due_date, paid, category_id FROM expenses WHERE user_id = $1  AND deleted = FALSE", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []postgres.Expense
	for rows.Next() {
		var e postgres.Expense
		if err := rows.Scan(&e.ID, &e.UserID, &e.Description, &e.Amount, &e.DueDate, &e.Paid, &e.CategoryID); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

// GetExpensesByUser retrieves all expenses for a specific user
func GetDeletedExpensesByUser(userID int) ([]postgres.Expense, error) {
	rows, err := postgres.DB.Query("SELECT id, user_id, description, amount, due_date, paid, category_id FROM expenses WHERE user_id = $1  AND deleted = TRUE", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []postgres.Expense
	for rows.Next() {
		var e postgres.Expense
		if err := rows.Scan(&e.ID, &e.UserID, &e.Description, &e.Amount, &e.DueDate, &e.Paid, &e.CategoryID); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

// GetExpense retrieves an expense by its ID
func GetExpense(expenseID, userID int) (*postgres.Expense, error) {
	var expense postgres.Expense
	err := postgres.DB.QueryRow("SELECT id, user_id, description, amount, due_date, paid, category_id FROM expenses WHERE id = $1 AND user_id = $2", expenseID, userID).Scan(&expense.ID, &expense.UserID, &expense.Description, &expense.Amount, &expense.DueDate, &expense.Paid, &expense.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &expense, nil
}

// CreateExpense inserts a new expense record into the database
func CreateExpense(userID int, description string, amount float64, dueDate string, categoryID sql.NullInt64) (int, error) {
	var expenseID int
	err := postgres.DB.QueryRow("INSERT INTO expenses (user_id, description, amount, due_date, category_id) VALUES ($1, $2, $3, $4, $5) RETURNING id", userID, description, amount, dueDate, categoryID).Scan(&expenseID)
	if err != nil {
		return 0, err
	}
	return expenseID, nil
}

// UpdateExpense modifies an existing expense record
func UpdateExpense(expenseID, userID int, description string, amount float64, dueDate string, paid bool, categoryID sql.NullInt64) error {
	_, err := postgres.DB.Exec("UPDATE expenses SET description = $1, amount = $2, due_date = $3, paid = $4, category_id = $5 WHERE id = $6 AND user_id = $7", description, amount, dueDate, paid, categoryID, expenseID, userID)
	return err
}

// DeleteExpense marks an expense as deleted
func DeleteExpense(expenseID, userID int) error {
	_, err := postgres.DB.Exec("UPDATE expenses SET deleted = TRUE WHERE id = $1 AND user_id = $2", expenseID, userID)
	return err
}

// RecoveryExpense marks an expense as not deleted
func RecoveryExpense(expenseID, userID int) error {
	_, err := postgres.DB.Exec("UPDATE expenses SET deleted = FALSE WHERE id = $1 AND user_id = $2", expenseID)
	return err
}
