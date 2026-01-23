package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // added postgres driver
	"github.com/mgmacri/pool-maintenance-app/internal/delivery"
	"github.com/mgmacri/pool-maintenance-app/internal/middleware"
	"github.com/mgmacri/pool-maintenance-app/internal/repository" // import repo
	"github.com/mgmacri/pool-maintenance-app/internal/usecase"    // import usecase
	"github.com/mgmacri/pool-maintenance-app/internal/version"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

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
	env := getEnvDefault("ENV", "dev")
	logLevelStr := strings.ToLower(getEnvDefault("LOG_LEVEL", "info"))
	lvl := parseLogLevel(logLevelStr)

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(lvl)
	logger, err := cfg.Build()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
	logger = logger.With(
		zap.String("service", "pool-maintenance-api"),
		zap.String("env", env),
		zap.String("version", version.Version),
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

	healthHandler := delivery.NewHealthHandler(logger) // no readiness checkers yet (placeholder)

	// 1. Initialize Repositories (Persistence)
	var poolRepo repository.PoolRepository
	dbHost := os.Getenv("DB_HOST")
	if dbHost != "" {
		dbUser := getEnvDefault("DB_USER", "pool_user")
		dbPass := getEnvDefault("DB_PASSWORD", "pool_pass")
		dbName := getEnvDefault("DB_NAME", "pool_db")
		// Assume standard port 5432 if not set
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPass, dbName)

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			logger.Fatal("Failed to open database connection", zap.Error(err))
		}

		// Verify connection
		if err := db.Ping(); err != nil {
			logger.Fatal("Failed to ping database", zap.Error(err))
		}

		poolRepo = repository.NewPostgresPoolRepository(db)
		logger.Info("Using Postgres Repository", zap.String("host", dbHost))
	} else {
		poolRepo = repository.NewInMemoryPoolRepo()
		logger.Info("Using In-Memory Repository (No persistence)")
	}

	// 2. Initialize UseCases (Business Logic)
	poolService := usecase.NewPoolUseCase(poolRepo)

	// 3. Initialize Handlers (Presentation)
	poolHandler := delivery.NewPoolHandler(logger, poolService)

	// Legacy single endpoint (backward compatibility)
	r.GET("/health", healthHandler.Check)
	// New split probes
	r.GET("/health/live", healthHandler.Live)
	r.GET("/health/ready", healthHandler.Ready)

	// 4. API Routes (Versioned)
	api := r.Group("/api/v1")
	{
		api.POST("/pools", poolHandler.CreatePool)
		api.GET("/pools/:id", poolHandler.GetPool)
		api.POST("/chemistry/analyze", poolHandler.AnalyzeChemistry)
	}

	logger.Info("starting server", zap.String("addr", ":8080"), zap.String("log_level", lvl.String()))
	if err := r.Run(":8080"); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}

func parseLogLevel(s string) zapcore.Level {
	switch s {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error", "err":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func getEnvDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
