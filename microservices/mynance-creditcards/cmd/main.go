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

	logger.Log.Info("[INFO] Starting MyFinance Categories...")

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
	// Grupo de rotas protegidas
	authRoutes := r.Group("/")
	authRoutes.Use(middleware.Auth())

	authRoutes.GET("/credit-cards", handlers.GetCreditCards)
	authRoutes.GET("/credit-cards/inactive", handlers.GetInactiveCreditCards)
	authRoutes.GET("/credit-cards/all", handlers.GetAllCreditCards)
	authRoutes.GET("/credit-cards/id", handlers.GetCreditCards)
	authRoutes.POST("/credit-cards", handlers.CreateCreditCard)
	authRoutes.PUT("/credit-cards", handlers.UpdateCreditCard)
	authRoutes.DELETE("/credit-cards", handlers.DeactivateCreditCard)
	authRoutes.POST("/credit-cards/activate", handlers.ActivateCreditCard)

	r.Run(":8080")
}
