package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvlerner/my-finance-api/internal/db"
	"github.com/jvlerner/my-finance-api/pkg/postgres"
)

// CreateCreditCardExpense handles credit card expense creation requests
func CreateCreditCardExpense(c *gin.Context) {
	var request postgres.CreditCardExpense
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	categoryID := sql.NullInt64{}
	if request.CategoryID != nil {
		categoryID.Valid = true
		categoryID.Int64 = int64(*request.CategoryID)
	}

	expenseID, err := db.CreateCreditCardExpense(request.CardID, userID, request.Description, request.Amount, request.PurchaseDate.Format("2006-01-02"), request.InstallmentCount, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create credit card expense"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": expenseID, "description": request.Description})
}

// GetCreditCardExpenses retrieves all active credit card expenses for a specific card
func GetCreditCardExpenses(c *gin.Context) {
	var request postgres.CreditCardExpense
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	expenses, err := db.GetCreditCardExpensesByCard(request.CardID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve credit card expenses"})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

// GetCreditCardExpense retrieves a credit card expense by ID
func GetCreditCardExpense(c *gin.Context) {
	var request postgres.CreditCardExpense
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	expense, err := db.GetCreditCardExpense(request.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve credit card expense"})
		return
	}
	if expense == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Credit card expense not found"})
		return
	}

	c.JSON(http.StatusOK, expense)
}

// UpdateCreditCardExpense modifies an existing credit card expense
func UpdateCreditCardExpense(c *gin.Context) {
	var request postgres.CreditCardExpense
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	categoryID := sql.NullInt64{}
	if request.CategoryID != nil {
		categoryID.Valid = true
		categoryID.Int64 = int64(*request.CategoryID)
	}

	err := db.UpdateCreditCardExpense(request.ID, userID, request.Description, request.Amount, request.PurchaseDate.Format("2006-01-02"), request.InstallmentCount, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update credit card expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credit card expense updated successfully"})
}

// DeleteCreditCardExpense marks a credit card expense as deleted
func DeleteCreditCardExpense(c *gin.Context) {
	var request postgres.CreditCardExpense
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	err := db.DeleteCreditCardExpense(request.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete credit card expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credit card expense deleted successfully"})
}

// RecoveryCreditCardExpense marks a credit card expense as not deleted
func RecoveryCreditCardExpense(c *gin.Context) {
	var request postgres.CreditCardExpense
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	err := db.RecoveryCreditCardExpense(request.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to recover credit card expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credit card expense recovered successfully"})
}
