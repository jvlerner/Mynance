package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvlerner/my-finance-api/internal/db"
	"github.com/jvlerner/my-finance-api/pkg/postgres"
)

// CreateExpense handles expense creation requests
func CreateExpense(c *gin.Context) {
	var request postgres.Expense
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	categoryID := sql.NullInt64{}
	if request.CategoryID != nil {
		categoryID.Valid = true
		categoryID.Int64 = int64(*request.CategoryID) // Converte int para int64
	}

	expenseID, err := db.CreateExpense(userID, request.Description, request.Amount, request.DueDate.Format("2006-01-02"), categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": expenseID, "description": request.Description})
}

// GetExpenses handles retrieving all active expenses for a user
func GetExpenses(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	expenses, err := db.GetExpensesByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expenses"})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

// GetExpense retrieves an expense by ID
func GetExpense(c *gin.Context) {
	var request postgres.Expense
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	expense, err := db.GetExpense(request.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expense"})
		return
	}
	if expense == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	c.JSON(http.StatusOK, expense)
}

// UpdateExpense modifies an existing expense
func UpdateExpense(c *gin.Context) {
	var request postgres.Expense
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	categoryID := sql.NullInt64{}
	if request.CategoryID != nil {
		categoryID.Valid = true
		categoryID.Int64 = int64(*request.CategoryID) // Converte int para int64
	}

	err := db.UpdateExpense(request.ID, userID, request.Description, request.Amount, request.DueDate.Format("2006-01-02"), request.Paid, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense updated successfully"})
}

// DeleteExpense marks an expense as deleted
func DeleteExpense(c *gin.Context) {
	var request postgres.Expense
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	err := db.DeleteExpense(request.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}

// RecoveryExpense marks an expense as not deleted
func RecoveryExpense(c *gin.Context) {
	var request postgres.Expense
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	err := db.RecoveryExpense(request.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to recover expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense recovered successfully"})
}
