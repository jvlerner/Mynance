package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBanks(c *gin.Context) {
	banks := []string{
		"Banco do Brasil",
		"Caixa Econômica Federal",
		"Itau",
		"Bradesco",
		"Santander",
		"Banco Safra",
		"BTG Pactual",
		"Banco Inter",
		"Nubank",
		"C6 Bank",
		"Banco Pan",
		"PagBank",
		"Neon",
		"Original",
		"Mercado Pago",
		"PicPay",
	}

	c.JSON(http.StatusOK, banks)
}
