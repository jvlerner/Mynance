package handlers

import (
	"errors"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jvlerner/my-finance-api/internal/db"
)

// ValidateToken verifies and returns the claims if the token is valid
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Pegar a data da última alteração de senha do banco de dados
	lastPasswordChange, err := db.UserLastPasswordChange(userDBName, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Se o token foi gerado antes da última alteração de senha, ele é inválido
	if claims.LastPasswordChange < lastPasswordChange.Unix() {
		return nil, errors.New("token invalid due to password change")
	}

	return claims, nil
}

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
