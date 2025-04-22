package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jvlerner/my-finance-api/pkg/postgres"

	"golang.org/x/crypto/bcrypt"
)

// UserExists checks if a user exists by email
func UserExists(email string) (bool, error) {
	var exists bool
	err := postgres.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// UserLastPasswordChange retrieves the last password change timestamp for a given user
func UserLastPasswordChange(userID int) (time.Time, error) {
	var lastPasswordChange time.Time
	err := postgres.DB.QueryRow("SELECT last_password_change FROM users WHERE id = $1", userID).Scan(&lastPasswordChange)
	if err != nil {
		if err == sql.ErrNoRows {
			// Retorna um timestamp zero se o usuário não for encontrado
			return time.Time{}, nil
		}
		return time.Time{}, err
	}
	return lastPasswordChange, nil
}

// CreateUser inserts a new user into the database
func CreateUser(name, email, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var userID int
	err = postgres.DB.QueryRow("INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id", name, email, string(hashedPassword)).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (*postgres.User, error) {
	var user postgres.User
	err := postgres.DB.QueryRow("SELECT id, name, email, password, active, created_at FROM users WHERE email = $1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Active, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetProfileByID retrieves user profile by ID
func GetProfileByID(userID int) (*postgres.Profile, error) {
	var user postgres.Profile
	err := postgres.DB.QueryRow("SELECT id, name, email, active, created_at FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Email, &user.Active, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user details
func UpdateUser(userID int, name string) error {
	_, err := postgres.DB.Exec("UPDATE users SET name = $1 WHERE id = $2", name, userID)
	return err
}

// UpdateUserPassword updates user password
func UpdateUserPassword(userID int, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = postgres.DB.Exec("UPDATE users SET password = $1, last_password_change = NOW() WHERE id = $2", string(hashedPassword), userID)
	return err
}

// DeleteUser marks a user as inactive instead of permanent deletion
func DeleteUser(userID int) error {
	_, err := postgres.DB.Exec("UPDATE users SET active = FALSE WHERE id = $1", userID)
	return err
}

// RecoverUser reactivates a user account
func RecoverUser(userID int) error {
	_, err := postgres.DB.Exec("UPDATE users SET active = TRUE WHERE id = $1", userID)
	return err
}
