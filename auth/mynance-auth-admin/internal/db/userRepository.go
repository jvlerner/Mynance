package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jvlerner/my-finance-api/pkg/postgres"
	"golang.org/x/crypto/bcrypt"
)

// UserExists checks if a user exists by email
func UserExists(dbName, email string) (bool, error) {
	db := postgres.GetDB(dbName)

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	return exists, err
}

// CreateServiceAccount inserts a new service user
func CreateServiceAccount(dbName, name, email, password string) (int, error) {
	db := postgres.GetDB(dbName)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var userID int
	err = db.QueryRow("INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id",
		name, email, string(hashedPassword), "service").Scan(&userID)
	return userID, err
}

// CreateUser inserts a new user
func CreateUser(dbName, name, email, password string) (int, error) {
	db := postgres.GetDB(dbName)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var userID int
	err = db.QueryRow("INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		name, email, string(hashedPassword)).Scan(&userID)
	return userID, err
}

func UserLastPasswordChange(dbName string, userID int) (time.Time, error) {
	db := postgres.GetDB(dbName)

	var lastPasswordChange time.Time
	err := db.QueryRow("SELECT last_password_change FROM users WHERE id = $1", userID).Scan(&lastPasswordChange)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}
	return lastPasswordChange, nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(dbName, email string) (*postgres.User, error) {
	db := postgres.GetDB(dbName)

	var user postgres.User
	err := db.QueryRow(`SELECT id, name, email, password, active, created_at, last_password_change, role 
		FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Active, &user.CreatedAt, &user.LastPasswordChange, &user.Role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetProfileByID retrieves user profile by ID
func GetProfileByID(dbName string, userID int) (*postgres.Profile, error) {
	db := postgres.GetDB(dbName)

	var user postgres.Profile
	err := db.QueryRow("SELECT id, name, email, active, created_at FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Active, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user details
func UpdateUser(dbName string, userID int, name string) error {
	db := postgres.GetDB(dbName)
	_, err := db.Exec("UPDATE users SET name = $1 WHERE id = $2", name, userID)
	return err
}

// UpdateUserPassword updates user password
func UpdateUserPassword(dbName string, userID int, password string) error {
	db := postgres.GetDB(dbName)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE users SET password = $1, last_password_change = NOW() WHERE id = $2",
		string(hashedPassword), userID)
	return err
}

// DeleteUser marks a user as inactive
func DeleteUser(dbName string, userID int) error {
	db := postgres.GetDB(dbName)
	_, err := db.Exec("UPDATE users SET active = FALSE WHERE id = $1", userID)
	return err
}

// RecoverUser reactivates a user account
func RecoverUser(dbName string, userID int) error {
	db := postgres.GetDB(dbName)
	_, err := db.Exec("UPDATE users SET active = TRUE WHERE id = $1", userID)
	return err
}
