package handlers

import (
	"net/http"

	"github.com/jvlerner/my-finance-api/internal/db"
	"github.com/jvlerner/my-finance-api/pkg/postgres"

	"github.com/gin-gonic/gin"
)

// CreateCreditCard handles credit card creation requests
func CreateCreditCard(c *gin.Context) {
	var request postgres.CreditCard
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	// Validação dos campos
	if request.Name == "" || request.Bank == "" ||
		request.LimitAmount < 0 || request.DueDay < 1 || request.DueDay > 31 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, invalid fields"})
		return
	}

	cardID, err := db.CreateCreditCard(userID, request.Name, request.Bank, request.LimitAmount, request.DueDay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create credit card"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": cardID, "name": request.Name})
}

// GetCreditCards handles retrieving all active credit cards
func GetCreditCards(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	cards, err := db.GetCreditCardsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve credit cards"})
		return
	}

	c.JSON(http.StatusOK, cards)
}

// GetCreditCards handles retrieving all active credit cards
func GetAllCreditCards(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	cards, err := db.GetAllCreditCardsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve credit cards"})
		return
	}

	c.JSON(http.StatusOK, cards)
}

func GetInactiveCreditCards(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	cards, err := db.GetInactiveCreditCardsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve inactive credit cards"})
		return
	}

	c.JSON(http.StatusOK, cards)
}

// GetCreditCard retrieves a credit card by ID
func GetCreditCard(c *gin.Context) {
	var request postgres.CreditCard
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	card, err := db.GetCreditCard(request.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve credit card"})
		return
	}
	if card == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Credit card not found"})
		return
	}

	c.JSON(http.StatusOK, card)
}

// UpdateCreditCard modifies an existing credit card
func UpdateCreditCard(c *gin.Context) {
	var request postgres.CreditCard
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	err := db.UpdateCreditCard(request.ID, userID, request.DueDay, request.Name, request.Bank, request.LimitAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update credit card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credit card updated successfully"})
}

// DeactivateCreditCard marks a credit card as inactive
func DeactivateCreditCard(c *gin.Context) {
	var request postgres.CreditCard
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	err := db.DeactivateCreditCard(request.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate credit card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credit card deactivated successfully"})
}

// ActivateCreditCard marks a credit card as active
func ActivateCreditCard(c *gin.Context) {
	var request postgres.CreditCard
	userID := c.MustGet("user_id").(int)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, unexpected fields"})
		return
	}

	err := db.ActivateCreditCard(request.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate credit card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credit card activated successfully"})
}
