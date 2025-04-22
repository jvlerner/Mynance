package db

import (
	"database/sql"

	"github.com/jvlerner/my-finance-api/pkg/postgres"
)

// CreatePayment inserts a new payment record into the database
func CreatePayment(expenseID, userID int, paidAt string, amount float64) (int, error) {
	var paymentID int
	err := postgres.DB.QueryRow("INSERT INTO payments (expense_id, user_id, paid_at, amount) VALUES ($1, $2, $3, $4) RETURNING id", expenseID, userID, paidAt, amount).Scan(&paymentID)
	if err != nil {
		return 0, err
	}
	return paymentID, nil
}

// GetPayment retrieves a payment by its ID
func GetPayment(paymentID, userID int) (*postgres.Payment, error) {
	var payment postgres.Payment
	err := postgres.DB.QueryRow("SELECT id, user_id, expense_id, paid_at, amount, deleted FROM payments WHERE id = $1 AND user_id = $2", paymentID, userID).Scan(&payment.ID, &payment.UserID, &payment.ExpenseID, &payment.PaidAt, &payment.Amount, &payment.Deleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

// GetPaymentsByExpense retrieves all payments for a specific expense
func GetPaymentsByExpense(expenseID, userID int) ([]postgres.Payment, error) {
	rows, err := postgres.DB.Query("SELECT id, user_id, expense_id, paid_at, amount FROM payments WHERE expense_id = $1 AND user_id = $2 AND deleted = FALSE", expenseID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []postgres.Payment
	for rows.Next() {
		var p postgres.Payment
		if err := rows.Scan(&p.ID, &p.UserID, &p.ExpenseID, &p.PaidAt, &p.Amount); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

// GetDeletedPaymentsByExpense retrieves all deleted payments for a specific expense
func GetDeletedPaymentsByExpense(expenseID, userID int) ([]postgres.Payment, error) {
	rows, err := postgres.DB.Query("SELECT id, user_id, expense_id, paid_at, amount FROM payments WHERE expense_id = $1 AND user_id = $2 AND deleted = TRUE", expenseID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []postgres.Payment
	for rows.Next() {
		var p postgres.Payment
		if err := rows.Scan(&p.ID, &p.UserID, &p.ExpenseID, &p.PaidAt, &p.Amount); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

// UpdatePayment modifies an existing payment record
func UpdatePayment(paymentID, userID int, paidAt string, amount float64) error {
	_, err := postgres.DB.Exec("UPDATE payments SET paid_at = $1, amount = $2 WHERE id = $3 AND user_id = $4", paidAt, amount, paymentID, userID)
	return err
}

// DeletePayment marks a payment as deleted
func DeletePayment(paymentID, userID int) error {
	_, err := postgres.DB.Exec("UPDATE payments SET deleted = TRUE WHERE id = $1 AND user_id = $2", paymentID, userID)
	return err
}

// RecoveryPayment marks a payment as not deleted
func RecoveryPayment(paymentID, userID int) error {
	_, err := postgres.DB.Exec("UPDATE payments SET deleted = FALSE WHERE id = $1 AND user_id = $2", paymentID, userID)
	return err
}
