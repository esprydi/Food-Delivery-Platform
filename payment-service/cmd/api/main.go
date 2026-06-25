package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"payment-service/config"
	"payment-service/internal/domain"
	handler "payment-service/internal/delivery/http"
	"payment-service/internal/infrastructure/rabbitmq"
	repoPostgres "payment-service/internal/repository/postgres"
	"payment-service/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/midtrans/midtrans-go"
	"gorm.io/driver/postgres" // standard for gorm
	"gorm.io/gorm"
)

func main() {
	// 1. Setup Structured Logging (slog)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	slog.Info("Starting Payment Service...")

	// 2. Load Config
	cfg := config.LoadConfig()

	// 3. Configure Midtrans
	midtrans.ServerKey = cfg.MidtransServerKey
	if cfg.IsProduction {
		midtrans.Environment = midtrans.Production
		slog.Info("Midtrans configured to Production mode")
	} else {
		midtrans.Environment = midtrans.Sandbox
		slog.Info("Midtrans configured to Sandbox mode")
	}

	// 4. Setup Database Connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// AutoMigrate
	err = db.AutoMigrate(&domain.Payment{})
	if err != nil {
		slog.Error("Failed to auto migrate database", "error", err)
	}

	// 5. Setup RabbitMQ Publisher
	eventPublisher, err := rabbitmq.NewRabbitMQPublisher(cfg.RabbitMQURL)
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ Publisher", "error", err)
	}

	// 6. Dependency Injection
	paymentRepo := repoPostgres.NewPaymentRepository(db)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo, eventPublisher)

	// 7. Setup RabbitMQ Consumer
	eventConsumer, err := rabbitmq.NewRabbitMQConsumer(cfg.RabbitMQURL, paymentUsecase)
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ Consumer", "error", err)
	} else {
		ctx := context.Background()
		err = eventConsumer.StartConsuming(ctx)
		if err != nil {
			slog.Error("Failed to start consuming", "error", err)
		}
	}

	// 8. Setup Echo Framework
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

	// Register Handlers
	handler.NewPaymentHandler(e, paymentUsecase)

	// 9. Start Server
	port := fmt.Sprintf(":%s", cfg.AppPort)
	slog.Info("Server listening", "port", port)
	e.Logger.Fatal(e.Start(port))
}
