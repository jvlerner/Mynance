package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvlerner/my-finance-api/internal/db"
	"github.com/jvlerner/my-finance-api/pkg/postgres"
)

// CreateCategory handles category creation requests
func CreateCategory(c *gin.Context) {
	var request postgres.Category
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	categoryID, err := db.CreateCategory(userID, request.Name, request.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": categoryID, "name": request.Name})
}

// GetCategories handles retrieving all active categories
func GetCategories(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	categories, err := db.GetCategories(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategories handles retrieving all active categories
func GetAllCategories(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	categories, err := db.GetAllCategories(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategories handles retrieving all active categories
func GetInactivateCategories(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	categories, err := db.GetInactiveCategories(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategory retrieves a category by ID
func GetCategory(c *gin.Context) {
	var request postgres.Category
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	category, err := db.GetCategory(userID, request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve category"})
		return
	}
	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategory modifies an existing category
func UpdateCategory(c *gin.Context) {
	var request postgres.Category
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	err := db.UpdateCategory(userID, request.ID, request.Name, request.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

// DeactivateCategory marks a category as inactive
func DeactivateCategory(c *gin.Context) {
	var request postgres.Category
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	err := db.DeactivateCategory(request.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deactivated successfully"})
}

// ActivateCategory marks a category as active
func ActivateCategory(c *gin.Context) {
	var request postgres.Category
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}
	err := db.ActivateCategory(request.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category activated successfully"})
}
