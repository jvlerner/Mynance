package db

import (
	"database/sql"

	"github.com/jvlerner/my-finance-api/pkg/postgres"
)

// CreateIncome inserts a new income record into the database
func CreateIncome(userID int, description string, amount float64, receivedAt string, isRecurring bool) (int, error) {
	var incomeID int
	err := postgres.DB.QueryRow("INSERT INTO incomes (user_id, description, amount, received_at, is_recurring) VALUES ($1, $2, $3, $4, $5) RETURNING id", userID, description, amount, receivedAt, isRecurring).Scan(&incomeID)
	if err != nil {
		return 0, err
	}
	return incomeID, nil
}

// GetIncome retrieves an income by its ID
func GetIncome(incomeID, userID int) (*postgres.Income, error) {
	var income postgres.Income
	err := postgres.DB.QueryRow("SELECT id, user_id, description, amount, received_at, is_recurring, deleted FROM incomes WHERE id = $1 AND user_id = $2", incomeID, userID).Scan(&income.ID, &income.UserID, &income.Description, &income.Amount, &income.ReceivedAt, &income.IsRecurring, &income.Deleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &income, nil
}

// GetIncomesByUser retrieves all incomes for a specific user
func GetIncomesByUser(userID int) ([]postgres.Income, error) {
	rows, err := postgres.DB.Query("SELECT id, user_id, description, amount, received_at, is_recurring FROM incomes WHERE user_id = $1 AND deleted = FALSE", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []postgres.Income
	for rows.Next() {
		var i postgres.Income
		if err := rows.Scan(&i.ID, &i.UserID, &i.Description, &i.Amount, &i.ReceivedAt, &i.IsRecurring); err != nil {
			return nil, err
		}
		incomes = append(incomes, i)
	}
	return incomes, nil
}

// GetDeletedIncomesByUser retrieves all deleted incomes for a specific user
func GetDeletedIncomesByUser(userID int) ([]postgres.Income, error) {
	rows, err := postgres.DB.Query("SELECT id, user_id, description, amount, received_at, is_recurring FROM incomes WHERE user_id = $1 AND deleted = TRUE", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []postgres.Income
	for rows.Next() {
		var i postgres.Income
		if err := rows.Scan(&i.ID, &i.UserID, &i.Description, &i.Amount, &i.ReceivedAt, &i.IsRecurring); err != nil {
			return nil, err
		}
		incomes = append(incomes, i)
	}
	return incomes, nil
}

// UpdateIncome modifies an existing income record
func UpdateIncome(incomeID, userID int, description string, amount float64, receivedAt string, isRecurring bool) error {
	_, err := postgres.DB.Exec("UPDATE incomes SET description = $1, amount = $2, received_at = $3, is_recurring = $4 WHERE id = $5 AND user_id = $6", description, amount, receivedAt, isRecurring, incomeID, userID)
	return err
}

// DeleteIncome marks an income as deleted
func DeleteIncome(incomeID, userID int) error {
	_, err := postgres.DB.Exec("UPDATE incomes SET deleted = TRUE WHERE id = $1 AND user_id = $2", incomeID, userID)
	return err
}

// RecoveryIncome marks an income as not deleted
func RecoveryIncome(incomeID, userID int) error {
	_, err := postgres.DB.Exec("UPDATE incomes SET deleted = FALSE WHERE id = $1 AND user_id = $2", incomeID, userID)
	return err
}
