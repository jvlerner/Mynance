package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvlerner/my-finance-api/pkg/auth"
)

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
		claims, err := auth.ValidateUserToken(tokenCookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - " + err.Error()})
			c.Abort()
			return
		}

		// Attach user info to the context
		c.Set("userId", claims.UserID)
		c.Set("userEmail", claims.Email)

		// Continue to next handler
		c.Next()
	}
}
