package main

import (
	"fmt"
	"log/slog"
	"os"

	"catalog-service/config"
	"catalog-service/internal/domain"
	handler "catalog-service/internal/delivery/http"
	repoPostgres "catalog-service/internal/repository/postgres"
	"catalog-service/internal/usecase"

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
	slog.Info("Starting Catalog Service...")

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

	// AutoMigrate (for development only)
	err = db.AutoMigrate(&domain.Restaurant{}, &domain.MenuItem{})
	if err != nil {
		slog.Error("Failed to auto migrate database", "error", err)
	}

	// 4. Dependency Injection
	restaurantRepo := repoPostgres.NewRestaurantRepository(db)
	menuRepo := repoPostgres.NewMenuRepository(db)
	catalogUsecase := usecase.NewCatalogUsecase(restaurantRepo, menuRepo)

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

	// Setup JWT Middleware configuration
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.JWTSecret),
	})

	// Register Handlers
	handler.NewCatalogHandler(e, catalogUsecase, jwtMiddleware)

	// 6. Start Server
	port := fmt.Sprintf(":%s", cfg.AppPort)
	slog.Info("Server listening", "port", port)
	e.Logger.Fatal(e.Start(port))
}
