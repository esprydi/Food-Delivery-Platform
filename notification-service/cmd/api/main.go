package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"notification-service/config"

	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		slog.Error("Failed to open a channel", "error", err)
		os.Exit(1)
	}
	defer ch.Close()

	// Ensure exchanges exist
	err = ch.ExchangeDeclare("order_exchange", "direct", true, false, false, false, nil)
	if err != nil {
		slog.Error("Failed to declare exchange", "error", err)
	}
	err = ch.ExchangeDeclare("payment_exchange", "direct", true, false, false, false, nil)
	if err != nil {
		slog.Error("Failed to declare exchange", "error", err)
	}

	// Setup Queues
	qOrder, _ := ch.QueueDeclare("notification_order_queue", true, false, false, false, nil)
	ch.QueueBind(qOrder.Name, "order.created", "order_exchange", false, nil)
	ch.QueueBind(qOrder.Name, "order.paid", "order_exchange", false, nil)

	qPayment, _ := ch.QueueDeclare("notification_payment_queue", true, false, false, false, nil)
	ch.QueueBind(qPayment.Name, "payment.success", "payment_exchange", false, nil)

	msgsOrder, _ := ch.Consume(qOrder.Name, "", false, false, false, false, nil)
	msgsPayment, _ := ch.Consume(qPayment.Name, "", false, false, false, false, nil)

	go func() {
		for d := range msgsOrder {
			var payload map[string]interface{}
			_ = json.Unmarshal(d.Body, &payload)
			slog.Info("🔔 NOTIFICATION [ORDER]: "+d.RoutingKey, "payload", payload)
			d.Ack(false)
		}
	}()

	go func() {
		for d := range msgsPayment {
			var payload map[string]interface{}
			_ = json.Unmarshal(d.Body, &payload)
			slog.Info("🔔 NOTIFICATION [PAYMENT]: "+d.RoutingKey, "payload", payload)
			d.Ack(false)
		}
	}()

	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	slog.Info("Notification Service starting on port " + cfg.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
