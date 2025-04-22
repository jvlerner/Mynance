package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvlerner/my-finance-api/internal/db"
	"github.com/jvlerner/my-finance-api/pkg/logger"
	"github.com/jvlerner/my-finance-api/pkg/postgres"
	"go.uber.org/zap"
)

func RegisterServiceAccount(c *gin.Context) {
	var user postgres.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	exists, err := db.UserExists(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email is already in use"})
		return
	}

	if !isValidPassword(user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long, include an uppercase letter, a lowercase letter, a number, and a special character."})
		return
	}

	userID, err := db.CreateServiceAccount(user.Name, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		logger.Log.Error("Failed to create user", zap.String("userName", user.Name), zap.String("userEmail", user.Email), zap.Error(err))
		return
	}

	logger.Log.Info("ServiceAccount registered successfully", zap.Int("userID", userID))
	c.JSON(http.StatusOK, gin.H{"message": "ServiceAccount registered successfully"})
}
