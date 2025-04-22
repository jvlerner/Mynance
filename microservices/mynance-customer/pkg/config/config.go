package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] Using system environment variables")
		return
	}
	log.Println("[INFO] Using .env variables")
}

// GetCORS returns a slice of allowed CORS origins
func GetCORS() []string {
	corsEnv := os.Getenv("CORS")
	if corsEnv != "" {
		// Permite múltiplos domínios separados por vírgula
		return strings.Split(corsEnv, ",")
	}
	// Default localhost para desenvolvimento
	return []string{"http://localhost:3000"}
}
