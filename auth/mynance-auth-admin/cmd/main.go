package main

import (
	"os"
	"sync"
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

var (
	wg          sync.WaitGroup
	userDBName  = os.Getenv("USER_DB_NAME")
	adminDBName = os.Getenv("ADMIN_DB_NAME")
)

func main() {
	config.LoadEnv()
	allowedOrigins := config.GetCORS()

	logger.InitLogger()
	defer logger.CloseLogger()

	logger.Log.Info("[INFO] Starting MyFinance Auth Admin/Service...")

	wg.Add(2)
	go func() {
		defer wg.Done()
		postgres.InitDB(userDBName,
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
	}()

	go func() {
		defer wg.Done()
		postgres.InitDB(adminDBName,
			os.Getenv("DB_ADMIN_USER"),
			os.Getenv("DB_ADMIN_PASSWORD"),
			os.Getenv("DB_ADMIN_HOST"),
			os.Getenv("DB_ADMIN_PORT"),
			os.Getenv("DB_ADMIN_NAME"),
		)
	}()
	wg.Wait()
	defer postgres.CloseAll()

	// Initialize the database metrics
	prometheus.SetDBForMonitoring(postgres.GetDB(userDBName), userDBName)
	defer prometheus.RemoveDBFromMonitoring(userDBName)

	prometheus.SetDBForMonitoring(postgres.GetDB(adminDBName), adminDBName)
	defer prometheus.RemoveDBFromMonitoring(adminDBName)

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

	r.POST("/auth/service/login", handlers.LoginService)
	r.POST("/auth/admin/login", handlers.LoginAdmin)

	seviceRoutes := r.Group("/auth/service")
	seviceRoutes.Use(middleware.ServiceAuth())
	seviceRoutes.GET("/validate-token", handlers.ValidateUserToken)

	adminRoutes := r.Group("/auth/admin")
	adminRoutes.Use(middleware.AdminAuth())
	adminRoutes.POST("/logout", handlers.LogoutAdmin)
	adminRoutes.POST("/register", handlers.RegisterServiceAccount)

	r.Run(":8080")
}
