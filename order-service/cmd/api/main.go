package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"order-service/config"
	"order-service/internal/domain"
	handler "order-service/internal/delivery/http"
	"order-service/internal/infrastructure/rabbitmq"
	repoPostgres "order-service/internal/repository/postgres"
	"order-service/internal/usecase"

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
	slog.Info("Starting Order Service...")

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

	// AutoMigrate
	err = db.AutoMigrate(&domain.Order{}, &domain.OrderItem{})
	if err != nil {
		slog.Error("Failed to auto migrate database", "error", err)
	}

	// 4. Setup RabbitMQ Publisher
	eventPublisher, err := rabbitmq.NewRabbitMQPublisher(cfg.RabbitMQURL)
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ", "error", err)
		// Usually we might exit here, but to allow local testing without RabbitMQ we can just warn
		// os.Exit(1) 
	} else {
		slog.Info("Connected to RabbitMQ successfully")
	}

	// 5. Dependency Injection
	orderRepo := repoPostgres.NewOrderRepository(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, eventPublisher)

	// 5.5 Setup RabbitMQ Consumer
	eventConsumer, err := rabbitmq.NewRabbitMQConsumer(cfg.RabbitMQURL, orderUsecase)
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ Consumer", "error", err)
	} else {
		// Start consumer in background
		ctx := context.Background()
		err = eventConsumer.StartConsuming(ctx)
		if err != nil {
			slog.Error("Failed to start consuming", "error", err)
		}
	}

	// 6. Setup Echo Framework
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			slog.Info("request", "uri", v.URI, "status", v.Status, "method", v.Method)
			return nil
		},
	}))

	// Setup JWT Middleware configuration
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.JWTSecret),
	})

	// Register Handlers
	handler.NewOrderHandler(e, orderUsecase, jwtMiddleware)

	// 7. Start Server
	port := fmt.Sprintf(":%s", cfg.AppPort)
	slog.Info("Server listening", "port", port)
	e.Logger.Fatal(e.Start(port))
}
