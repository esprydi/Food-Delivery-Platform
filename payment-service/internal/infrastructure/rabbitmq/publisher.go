package rabbitmq

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"payment-service/internal/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPublisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQPublisher(url string) (domain.PaymentEventPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &EventPublisher{
		conn: conn,
		ch:   ch,
	}, nil
}

func (p *EventPublisher) PublishPaymentSuccess(ctx context.Context, orderID string) error {
	payload := map[string]interface{}{
		"event_id":   orderID + "_paid",
		"event_type": "PAYMENT_SUCCESS",
		"timestamp":  time.Now().Format(time.RFC3339),
		"payload": map[string]interface{}{
			"order_id": orderID,
			"status":   "SUCCESS",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = p.ch.PublishWithContext(ctx,
		"food_delivery_events", // exchange
		"payment.success",      // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		slog.Error("Failed to publish PaymentSuccess event", "error", err)
		return err
	}

	slog.Info("Successfully published PaymentSuccess event", "order_id", orderID)
	return nil
}
