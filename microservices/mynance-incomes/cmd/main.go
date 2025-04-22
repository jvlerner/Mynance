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

	logger.InitLogger()
	defer logger.CloseLogger()

	logger.Log.Info("[INFO] Starting MyFinance API...")

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
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
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

	// Grupo de rotas protegidas
	authRoutes := r.Group("/")
	authRoutes.Use(middleware.Auth())
	authRoutes.GET("/user/me", handlers.GetUserProfile)
	authRoutes.POST("/user/password", handlers.UpdateUserPassword)
	authRoutes.POST("/user/name", handlers.UpdateUserName)
	authRoutes.DELETE("/user", handlers.DeactivateUser)
	authRoutes.POST("/user/activate", handlers.ActivateUser)

	authRoutes.GET("/expenses", handlers.GetExpenses)
	// authRoutes.GET("/expenses/inactive", handlers.GetExpenses)
	// authRoutes.GET("/expenses/all", handlers.GetExpenses)
	authRoutes.GET("/expenses/id", handlers.GetExpense)
	authRoutes.POST("/expenses", handlers.CreateExpense)
	authRoutes.PUT("/expenses", handlers.UpdateExpense)
	authRoutes.DELETE("/expenses", handlers.DeleteExpense)
	authRoutes.POST("/expenses/activate", handlers.RecoveryExpense)

	authRoutes.GET("/incomes", handlers.GetIncomes)
	// authRoutes.GET("/incomes/inactive", handlers.GetInactiveIncomes)
	// authRoutes.GET("/incomes/all", handlers.GetAllIncomes)
	authRoutes.GET("/incomes/id", handlers.GetIncome)
	authRoutes.POST("/incomes", handlers.CreateIncome)
	authRoutes.PUT("/incomes", handlers.UpdateIncome)
	authRoutes.DELETE("/incomes", handlers.DeleteIncome)
	authRoutes.POST("/incomes/activate", handlers.RecoveryIncome)

	authRoutes.GET("/credit-cards", handlers.GetCreditCards)
	authRoutes.GET("/credit-cards/inactive", handlers.GetInactiveCreditCards)
	authRoutes.GET("/credit-cards/all", handlers.GetAllCreditCards)
	authRoutes.GET("/credit-cards/id", handlers.GetCreditCards)
	authRoutes.POST("/credit-cards", handlers.CreateCreditCard)
	authRoutes.PUT("/credit-cards", handlers.UpdateCreditCard)
	authRoutes.DELETE("/credit-cards", handlers.DeactivateCreditCard)
	authRoutes.POST("/credit-cards/activate", handlers.ActivateCreditCard)

	authRoutes.GET("/credit-cards/expenses", handlers.GetCreditCardExpenses)
	// authRoutes.GET("/credit-cards/inactive", handlers.GetInactiveCreditCardExpenses)
	// authRoutes.GET("/credit-cards/all", handlers.GetAllCreditCardExpenses)
	authRoutes.GET("/credit-cards/expenses/id", handlers.GetCreditCardExpense)
	authRoutes.POST("/credit-cards/expenses", handlers.CreateCreditCardExpense)
	authRoutes.PUT("/credit-cards/expenses", handlers.UpdateCreditCardExpense)
	authRoutes.DELETE("/credit-cards/expenses", handlers.DeleteCreditCardExpense)
	authRoutes.POST("/credit-cards/expenses/activate", handlers.RecoveryCreditCardExpense)

	authRoutes.GET("/categories", handlers.GetCategories)
	authRoutes.GET("/categories/inactive", handlers.GetInactivateCategories)
	authRoutes.GET("/categories/all", handlers.GetAllCategories)
	authRoutes.GET("/categories/id", handlers.GetCategory)
	authRoutes.POST("/categories", handlers.CreateCategory)
	authRoutes.PUT("/categories", handlers.UpdateCategory)
	authRoutes.DELETE("/categories", handlers.DeactivateCategory)
	authRoutes.POST("/categories/activate", handlers.ActivateCategory)

	authRoutes.GET("/banks", handlers.GetBanks)
	// Iniciar o servidor
	r.Run(":8080")
}
