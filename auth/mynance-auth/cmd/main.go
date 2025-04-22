package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/jvlerner/my-finance-api/internal/handlers"
	"github.com/jvlerner/my-finance-api/internal/middleware"
	"github.com/jvlerner/my-finance-api/pkg/config"
	"github.com/jvlerner/my-finance-api/pkg/logger"
	"github.com/jvlerner/my-finance-api/pkg/postgres"
	"github.com/jvlerner/my-finance-api/pkg/prometheus"
)

func main() {
	config.LoadEnv()
	allowedOrigins := config.GetCORS()

	logger.InitLogger()
	defer logger.CloseLogger()

	logger.Log.Info("[INFO] Starting MyFinance Auth...")

	postgres.InitDB()
	defer postgres.CloseDB()

	// Initialize the database metrics
	prometheus.SetDBForMonitoring(postgres.DB)
	defer prometheus.CloseDBForMonitoring()

	// Initialize the metrics
	prometheus.Init()
	defer prometheus.Close()

	r := gin.Default()
	middleware.StartRateLimiter()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-User-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.Prometheus())
	r.Use(middleware.RateLimit())

	// Rotas de autenticação
	r.POST("/auth/register", handlers.Register)
	r.POST("/auth/login", handlers.Login)
	r.POST("/auth/logout", handlers.Logout)

	// Iniciar o servidor
	r.Run(":8080")
}
