package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateUserToken(c *gin.Context) {
	// Tenta obter o token do cookie
	token := c.Request.Header.Get("X-User-Token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid": false,
			"error": "Invalid or expired token",
		})
		return
	}

	// Valida o token
	claims, err := ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"valid": false,
			"error": "Invalid or expired token",
		})
		return
	}

	// Retorna os dados do usuário do token
	c.JSON(http.StatusOK, gin.H{
		"valid":     true,
		"userId":    claims.UserID,
		"email":     claims.Email,
		"role":      claims.Role,
		"expiresAt": claims.ExpiresAt,
	})
}
