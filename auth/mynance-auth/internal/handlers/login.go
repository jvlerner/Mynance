package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jvlerner/my-finance-api/internal/db"
	"github.com/jvlerner/my-finance-api/pkg/logger"
	"github.com/jvlerner/my-finance-api/pkg/postgres"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var user postgres.User
	var tokenExpiration time.Duration = 12

	// Verifica se há um token válido no cookie
	tokenCookie, err := c.Request.Cookie("token")
	if err == nil {
		// Valida o token existente
		claims, err := ValidateToken(tokenCookie.Value)
		if err == nil {
			logger.Log.Info("User already has a valid token", zap.Int("userID", claims.UserID))

			// Atualiza o token e o cookie
			newToken, err := GenerateToken(claims.UserID, claims.Email, claims.Role)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
				return
			}

			// Define o novo token no cookie
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "token",
				Value:    newToken,
				Expires:  time.Now().Add(tokenExpiration * time.Hour),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				Path:     "/",
			})

			c.JSON(http.StatusOK, gin.H{"message": "Token refreshed successfully"})
			return
		}
	}

	// Se não houver token válido, proceder com o login normal
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	exists, err := db.UserExists(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	storedUser, err := db.GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := GenerateToken(storedUser.ID, storedUser.Email, storedUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(tokenExpiration * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	logger.Log.Info("User logged in", zap.Int("userID", storedUser.ID))
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
