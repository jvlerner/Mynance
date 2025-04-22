package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LogoutAdmin(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token-admin",
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
