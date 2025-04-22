package db

import (
	"database/sql"

	"github.com/jvlerner/my-finance-api/pkg/postgres"
)

// CreateCreditCardExpense inserts a new credit card expense record into the database
func CreateCreditCardExpense(cardID, userID int, description string, amount float64, purchaseDate string, installmentCount int, categoryID sql.NullInt64) (int, error) {
	var expenseID int
	err := postgres.DB.QueryRow("INSERT INTO credit_card_expenses (card_id, user_id, description, amount, purchase_date, installment_count, category_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", cardID, userID, description, amount, purchaseDate, installmentCount, categoryID).Scan(&expenseID)
	if err != nil {
		return 0, err
	}
	return expenseID, nil
}

// GetCreditCardExpense retrieves a credit card expense by its ID
func GetCreditCardExpense(expenseID, userID int) (*postgres.CreditCardExpense, error) {
	var expense postgres.CreditCardExpense
	err := postgres.DB.QueryRow("SELECT id, user_id, card_id, description, amount, purchase_date, installment_count, category_id FROM credit_card_expenses WHERE id = $1 AND user_id = $2", expenseID, userID).Scan(&expense.ID, &expense.UserID, &expense.CardID, &expense.Description, &expense.Amount, &expense.PurchaseDate, &expense.InstallmentCount, &expense.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &expense, nil
}

// GetCreditCardExpensesByCard retrieves all expenses for a specific credit card
func GetCreditCardExpensesByCard(cardID, userID int) ([]postgres.CreditCardExpense, error) {
	rows, err := postgres.DB.Query("SELECT id, user_id, description, amount, purchase_date, installment_count, category_id FROM credit_card_expenses WHERE card_id = $1 AND user_id = $2 AND deleted = FALSE", cardID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []postgres.CreditCardExpense
	for rows.Next() {
		var e postgres.CreditCardExpense
		if err := rows.Scan(&e.ID, &e.UserID, &e.Description, &e.Amount, &e.PurchaseDate, &e.InstallmentCount, &e.CategoryID); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

// GetDeletedCreditCardExpensesByCard retrieves all deleted expenses for a specific credit card
func GetDeletedCreditCardExpensesByCard(cardID, userID int) ([]postgres.CreditCardExpense, error) {
	rows, err := postgres.DB.Query("SELECT id, user_id, description, amount, purchase_date, installment_count, category_id FROM credit_card_expenses WHERE card_id = $1 AND user_id = $2 AND deleted = TRUE", cardID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []postgres.CreditCardExpense
	for rows.Next() {
		var e postgres.CreditCardExpense
		if err := rows.Scan(&e.ID, &e.UserID, &e.Description, &e.Amount, &e.PurchaseDate, &e.InstallmentCount, &e.CategoryID); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

// UpdateCreditCardExpense modifies an existing credit card expense record
func UpdateCreditCardExpense(expenseID, userID int, description string, amount float64, purchaseDate string, installmentCount int, categoryID sql.NullInt64) error {
	_, err := postgres.DB.Exec("UPDATE credit_card_expenses SET description = $1, amount = $2, purchase_date = $3, installment_count = $4, category_id = $5 WHERE id = $6 AND user_id = $7", description, amount, purchaseDate, installmentCount, categoryID, expenseID, userID)
	return err
}

// DeleteCreditCardExpense marks a credit card expense as deleted
func DeleteCreditCardExpense(expenseID, userID int) error {
	_, err := postgres.DB.Exec("UPDATE credit_card_expenses SET deleted = TRUE WHERE id = $1 AND user_id = $2", expenseID, userID)
	return err
}

// RecoveryCreditCardExpense marks a credit card expense as not deleted
func RecoveryCreditCardExpense(expenseID, userID int) error {
	_, err := postgres.DB.Exec("UPDATE credit_card_expenses SET deleted = FALSE WHERE id = $1 AND user_id = $2", expenseID, userID)
	return err
}
