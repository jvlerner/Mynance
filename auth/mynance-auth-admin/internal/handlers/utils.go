package handlers

import (
	"os"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jvlerner/my-finance-api/internal/db"
)

var (
	userDBName  = os.Getenv("USER_DB_NAME")
	adminDBName = os.Getenv("ADMIN_DB_NAME")
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
	lastPasswordChange, err := db.UserLastPasswordChange(adminDBName, userID)
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
