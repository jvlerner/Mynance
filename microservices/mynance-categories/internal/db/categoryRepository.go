package db

import (
	"database/sql"

	"github.com/jvlerner/my-finance-api/pkg/postgres"
)

// CreateCategory inserts a new category record into the database
func CreateCategory(userID int, name string, color string) (int, error) {
	var categoryID int
	err := postgres.DB.QueryRow("INSERT INTO categories (name, color, user_id) VALUES ($1, $2, $3) RETURNING id", name, color, userID).Scan(&categoryID)
	if err != nil {
		return 0, err
	}
	return categoryID, nil
}

// GetCategory retrieves a category by its ID
func GetCategory(userID, categoryID int) (*postgres.Category, error) {
	var category postgres.Category
	err := postgres.DB.QueryRow("SELECT id, user_id, name, active, created_at FROM categories WHERE id = $1 AND user_id = $2", categoryID, userID).Scan(&category.ID, &category.UserID, &category.Name, &category.Active, &category.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// GetCategories retrieves all active categories
func GetCategories(userID int) ([]postgres.Category, error) {
	rows, err := postgres.DB.Query("SELECT id, user_id, name, color, active, created_at FROM categories WHERE active = TRUE AND user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []postgres.Category
	for rows.Next() {
		var c postgres.Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Color, &c.Active, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// GetCategories retrieves all active categories
func GetAllCategories(userID int) ([]postgres.Category, error) {
	rows, err := postgres.DB.Query("SELECT id, user_id, name, color, active, created_at FROM categories WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []postgres.Category
	for rows.Next() {
		var c postgres.Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Color, &c.Active, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// GetInactiveCategories retrieves all inactive categories
func GetInactiveCategories(userID int) ([]postgres.Category, error) {
	rows, err := postgres.DB.Query("SELECT id, user_id, name, color, active, created_at FROM categories WHERE active = FALSE AND user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []postgres.Category
	for rows.Next() {
		var c postgres.Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Color, &c.Active, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// UpdateCategory modifies an existing category record
func UpdateCategory(userID, categoryID int, name, color string) error {
	_, err := postgres.DB.Exec("UPDATE categories SET name = $1, color = $2 WHERE id = $3 AND user_id= $4", name, color, categoryID, userID)
	return err
}

// DeactivateCategory marks a category as inactive
func DeactivateCategory(userID, categoryID int) error {
	_, err := postgres.DB.Exec("UPDATE categories SET active = FALSE WHERE id = $1 AND user_id= $2", categoryID, userID)
	return err
}

// ActivateCategory marks a category as active
func ActivateCategory(userID, categoryID int) error {
	_, err := postgres.DB.Exec("UPDATE categories SET active = TRUE WHERE id = $1 AND user_id= $2", categoryID, userID)
	return err
}
