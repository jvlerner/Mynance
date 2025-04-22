package db

import (
	"database/sql"

	"github.com/jvlerner/my-finance-api/pkg/postgres"
)

// CreateCreditCard inserts a new credit card record into the database
func CreateCreditCard(userID int, name string, bank string, limitAmount float64, dueDay int) (int, error) {
	var cardID int
	err := postgres.DB.QueryRow("INSERT INTO credit_cards (user_id, name, bank, limit_amount, due_day) VALUES ($1, $2, $3, $4, $5) RETURNING id", userID, name, bank, limitAmount, dueDay).Scan(&cardID)
	if err != nil {
		return 0, err
	}
	return cardID, nil
}

// GetCreditCard retrieves a credit card by its ID
func GetCreditCard(cardID, userID int) (*postgres.CreditCard, error) {
	var card postgres.CreditCard
	err := postgres.DB.QueryRow("SELECT id, user_id, name, bank, limit_amount, due_day, active FROM credit_cards WHERE id = $1 AND user_id = $2", cardID, userID).Scan(&card.ID, &card.UserID, &card.Name, &card.Bank, &card.LimitAmount, &card.DueDay, &card.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &card, nil
}

// GetCreditCardsByUser retrieves all active credit cards for a specific user
func GetCreditCardsByUser(userID int) ([]postgres.CreditCard, error) {
	rows, err := postgres.DB.Query("SELECT id, user_id, name, bank, limit_amount, due_day FROM credit_cards WHERE user_id = $1 AND active = TRUE", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []postgres.CreditCard
	for rows.Next() {
		var c postgres.CreditCard
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Bank, &c.LimitAmount, &c.DueDay); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, nil
}

// GetCreditCardsByUser retrieves all active credit cards for a specific user
func GetAllCreditCardsByUser(userID int) ([]postgres.CreditCard, error) {
	rows, err := postgres.DB.Query("SELECT id, user_id, name, bank, limit_amount, due_day, active FROM credit_cards WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []postgres.CreditCard
	for rows.Next() {
		var c postgres.CreditCard
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Bank, &c.LimitAmount, &c.DueDay, &c.Active); err != nil {
			return nil, err
		}
		cards = append(cards, c)

	}

	return cards, nil
}

// GetInactiveCreditCardsByUser retrieves all inactive credit cards for a specific user
func GetInactiveCreditCardsByUser(userID int) ([]postgres.CreditCard, error) {
	rows, err := postgres.DB.Query("SELECT id, user_id, name, bank, limit_amount, due_day FROM credit_cards WHERE user_id = $1 AND active = FALSE", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []postgres.CreditCard
	for rows.Next() {
		var c postgres.CreditCard
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Bank, &c.LimitAmount, &c.DueDay); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, nil
}

// UpdateCreditCard modifies an existing credit card record
func UpdateCreditCard(cardID, userID, dueDay int, name string, bank string, limitAmount float64) error {
	_, err := postgres.DB.Exec("UPDATE credit_cards SET name = $1, bank = $2, limit_amount = $3, due_day = $4 WHERE id = $5 AND user_id = $6", name, bank, limitAmount, dueDay, cardID, userID)
	return err
}

// DeactivateCreditCard marks a credit card as inactive
func DeactivateCreditCard(cardID, userID int) error {
	_, err := postgres.DB.Exec("UPDATE credit_cards SET active = FALSE WHERE id = $1 AND user_id = $2", cardID, userID)
	return err
}

// ActivateCreditCard marks a credit card as active
func ActivateCreditCard(cardID, userID int) error {
	_, err := postgres.DB.Exec("UPDATE credit_cards SET active = TRUE WHERE id = $1 AND user_id = $2", cardID, userID)
	return err
}
