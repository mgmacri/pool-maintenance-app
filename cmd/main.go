package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mgmacri/pool-maintenance-app/internal/delivery"
	"github.com/mgmacri/pool-maintenance-app/internal/middleware"
	"github.com/mgmacri/pool-maintenance-app/internal/version"
	"go.uber.org/zap"

	// Swagger imports
	_ "github.com/mgmacri/pool-maintenance-app/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Pool Maintenance API
// @version         1.0
// @description     This is a sample API for pool maintenance, demonstrating Clean Architecture, structured logging, and DevOps best practices.
// @termsOfService  https://github.com/mgmacri/pool-maintenance-app
// @contact.name    Matthew G. Macri
// @contact.url     https://github.com/mgmacri
// @contact.email   mgmacri@example.com
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
// @host            localhost:8080
// @BasePath        /
// @schemes         http
func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	logger = logger.With(
		zap.String("service", "pool-maintenance-api"),
		zap.String("env", env),
		zap.String("version", version.Version),
		zap.String("trace_id", ""), // Placeholder for future tracing
	)
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("failed to sync logger", zap.Error(err))
		}
	}()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.ZapLogger(logger))

	// Register Swagger UI route after router is initialized
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	healthHandler := delivery.NewHealthHandler(logger)
	r.GET("/health", healthHandler.Check)

	logger.Info("starting server", zap.String("addr", ":8080"))
	if err := r.Run(":8080"); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
