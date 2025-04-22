package middleware

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jvlerner/my-finance-api/internal/db"
)

type Claims struct {
	UserID             int    `json:"user_id"`
	Email              string `json:"email"`
	LastPasswordChange int64  `json:"last_password_change"`
	jwt.StandardClaims
}

func GenerateToken(userID int, email string) (string, error) {
	lastPasswordChange, err := db.UserLastPasswordChange(userID)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id":              userID,
		"email":                email,
		"last_password_change": lastPasswordChange.Unix(),
		"exp":                  time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

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
	lastPasswordChange, err := db.UserLastPasswordChange(claims.UserID)
	if err != nil {
		return nil, err
	}

	// Se o token foi gerado antes da última alteração de senha, ele é inválido
	if claims.LastPasswordChange < lastPasswordChange.Unix() {
		return nil, errors.New("token invalid due to password change")
	}

	return claims, nil
}

// Auth extracts and validates the JWT token from the cookie
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the cookie
		tokenCookie, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - No token provided"})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := ValidateToken(tokenCookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - " + err.Error()})
			c.Abort()
			return
		}

		// Attach user info to the context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		// Continue to next handler
		c.Next()
	}
}
