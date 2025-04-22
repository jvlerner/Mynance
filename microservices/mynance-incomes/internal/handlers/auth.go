package handlers

import (
	"net/http"
	"regexp"
	"time"

	"github.com/jvlerner/my-finance-api/internal/db"
	"github.com/jvlerner/my-finance-api/internal/middleware"
	"github.com/jvlerner/my-finance-api/pkg/logger"
	"github.com/jvlerner/my-finance-api/pkg/postgres"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func isValidPassword(password string) bool {
	var (
		upperCase = regexp.MustCompile(`[A-Z]`)
		lowerCase = regexp.MustCompile(`[a-z]`)
		digit     = regexp.MustCompile(`\d`)
		special   = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	)
	return len(password) >= 8 && upperCase.MatchString(password) && lowerCase.MatchString(password) && digit.MatchString(password) && special.MatchString(password)
}

func Register(c *gin.Context) {
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

	userID, err := db.CreateUser(user.Name, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	logger.Log.Info("User registered successfully", zap.Int("userID", userID))
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var user postgres.User

	// Verifica se há um token válido no cookie
	tokenCookie, err := c.Request.Cookie("token")
	if err == nil {
		// Valida o token existente
		claims, err := middleware.ValidateToken(tokenCookie.Value)
		if err == nil {
			logger.Log.Info("User already has a valid token", zap.Int("userID", claims.UserID))

			// Atualiza o token e o cookie
			newToken, err := middleware.GenerateToken(claims.UserID, claims.Email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
				return
			}

			// Define o novo token no cookie
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "token",
				Value:    newToken,
				Expires:  time.Now().Add(24 * time.Hour),
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

	// Gerar novo token
	token, err := middleware.GenerateToken(storedUser.ID, storedUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Criar cookie do token
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	logger.Log.Info("User logged in", zap.Int("userID", storedUser.ID))
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Logout efetuado com sucesso"})
}
