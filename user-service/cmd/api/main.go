package main

import (
	"fmt"
	"log/slog"
	"os"

	"user-service/config"
	"user-service/internal/domain"
	handler "user-service/internal/delivery/http"
	repoPostgres "user-service/internal/repository/postgres"
	"user-service/internal/usecase"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres" // standard for gorm
	"gorm.io/gorm"
)

func main() {
	// 1. Setup Structured Logging (slog)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	slog.Info("Starting User Service...")

	// 2. Load Config
	cfg := config.LoadConfig()

	// 3. Setup Database Connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// AutoMigrate (for development only, use migrations in prod)
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		slog.Error("Failed to auto migrate database", "error", err)
	}

	// 4. Dependency Injection
	userRepo := repoPostgres.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, cfg.JWTSecret)

	// 5. Setup Echo Framework
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	
	// Logger middleware using slog
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			slog.Info("request",
				"uri", v.URI,
				"status", v.Status,
				"method", v.Method,
			)
			return nil
		},
	}))

	// Register Handlers
	handler.NewUserHandler(e, userUsecase)

	// Setup JWT Middleware for protected routes
	// Note: We need to apply this selectively to the handler later
	e.Group("/api/v1/users").Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.JWTSecret),
	}))

	// 6. Start Server
	port := fmt.Sprintf(":%s", cfg.AppPort)
	slog.Info("Server listening", "port", port)
	e.Logger.Fatal(e.Start(port))
}
