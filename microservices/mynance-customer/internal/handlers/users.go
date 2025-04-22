package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvlerner/my-finance-api/internal/db"
	"github.com/jvlerner/my-finance-api/pkg/logger"
	"go.uber.org/zap"
)

// GetUserProfile retrieves the profile of the logged-in user
func GetUserProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	userEmail := c.MustGet("user_email").(string)

	exists, err := db.UserExists(userEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	profile, err := db.GetProfileByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	logger.Log.Info("User profile retrieved", zap.Int("userID", userID))
	c.JSON(http.StatusOK, gin.H{"user": profile})
}

// UpdateUserName updates the name of the logged-in user
func UpdateUserName(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field name not provided"})
		return
	}

	if err := db.UpdateUser(userID, request.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update name"})
		return
	}

	logger.Log.Info("User name updated", zap.Int("userID", userID))
	c.JSON(http.StatusOK, gin.H{"message": "Name updated successfully"})
}

// UpdateUserPassword updates the password of the logged-in user
func UpdateUserPassword(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	var request struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field password not provided"})
		return
	}

	if !isValidPassword(request.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long, include an uppercase letter, a lowercase letter, a number, and a special character."})
		return
	}

	if err := db.UpdateUserPassword(userID, request.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	logger.Log.Info("User password updated", zap.Int("userID", userID))
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// DeactivateUser marks the user account as inactive
func DeactivateUser(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	if err := db.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate account"})
		return
	}

	logger.Log.Info("User account deactivated", zap.Int("userID", userID))
	c.JSON(http.StatusOK, gin.H{"message": "Account deactivated successfully"})
}

// ActivateUser reactivates the user account
func ActivateUser(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	if err := db.RecoverUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate account"})
		return
	}

	logger.Log.Info("User account activated", zap.Int("userID", userID))
	c.JSON(http.StatusOK, gin.H{"message": "Account activated successfully"})
}
