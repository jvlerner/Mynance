package handlers

import (
	"errors"
	"os"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jvlerner/my-finance-api/internal/db"
)

type Claims struct {
	UserID             int    `json:"userId"`
	Email              string `json:"email"`
	Role               string `json:"role"`
	LastPasswordChange int64  `json:"lastPasswordChange"`
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

func GenerateToken(userID int, email, role string) (string, int64, error) {
	lastPasswordChange, err := db.UserLastPasswordChange(userID)
	if err != nil {
		return "", 0, err
	}

	exp := time.Now().Add(24 * time.Hour).Unix()
	if role == "service" {
		exp = time.Now().Add(90 * 24 * time.Hour).Unix()
	}

	claims := jwt.MapClaims{
		"userId":             userID,
		"email":              email,
		"role":               role,
		"lastPasswordChange": lastPasswordChange.Unix(),
		"exp":                exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	assignedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", 0, err
	}
	return assignedToken, exp, nil
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
