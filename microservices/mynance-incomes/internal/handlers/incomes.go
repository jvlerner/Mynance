package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvlerner/my-finance-api/internal/db"
	"github.com/jvlerner/my-finance-api/pkg/postgres"
)

// CreateIncome handles income creation requests
func CreateIncome(c *gin.Context) {
	var request postgres.Income
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	incomeID, err := db.CreateIncome(userID, request.Description, request.Amount, request.ReceivedAt.Format("2006-01-02"), request.IsRecurring)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create income"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": incomeID, "description": request.Description})
}

// GetIncomes handles retrieving all active incomes for a user
func GetIncomes(c *gin.Context) {
	var request postgres.Income
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	incomes, err := db.GetIncomesByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incomes"})
		return
	}

	c.JSON(http.StatusOK, incomes)
}

// GetIncome retrieves an income by ID
func GetIncome(c *gin.Context) {
	var request postgres.Income
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	income, err := db.GetIncome(request.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve income"})
		return
	}
	if income == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Income not found"})
		return
	}

	c.JSON(http.StatusOK, income)
}

// UpdateIncome modifies an existing income
func UpdateIncome(c *gin.Context) {
	var request postgres.Income
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	err := db.UpdateIncome(request.ID, userID, request.Description, request.Amount, request.ReceivedAt.Format("2006-01-02"), request.IsRecurring)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update income"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Income updated successfully"})
}

// DeleteIncome marks an income as deleted
func DeleteIncome(c *gin.Context) {
	var request postgres.Income
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	err := db.DeleteIncome(request.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete income"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Income deleted successfully"})
}

// DeleteIncome marks an income as not deleted
func RecoveryIncome(c *gin.Context) {
	var request postgres.Income
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	err := db.RecoveryIncome(request.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to recover income"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Income recovered successfully"})
}
