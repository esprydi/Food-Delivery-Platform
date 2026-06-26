package rabbitmq

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"order-service/internal/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMQPublisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQPublisher(url string) (domain.EventPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Ensure the exchange exists
	err = ch.ExchangeDeclare(
		"food_delivery_events", // name
		"topic",                // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return nil, err
	}

	return &rabbitMQPublisher{
		conn: conn,
		ch:   ch,
	}, nil
}

func (p *rabbitMQPublisher) PublishOrderCreated(ctx context.Context, order *domain.Order) error {
	// Create the event payload exactly as defined in PRD
	payload := map[string]interface{}{
		"event_id":   order.ID + "_created",
		"event_type": "ORDER_CREATED",
		"timestamp":  time.Now().Format(time.RFC3339),
		"payload": map[string]interface{}{
			"order_id":       order.ID,
			"customer_id":    order.CustomerID,
			"customer_email": order.CustomerEmail,
			"restaurant_id":  order.RestaurantID,
			"total_amount":   order.TotalAmount,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = p.ch.PublishWithContext(ctx,
		"food_delivery_events", // exchange
		"order.created",        // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		slog.Error("Failed to publish OrderCreated event", "error", err)
		return err
	}

	slog.Info("Successfully published OrderCreated event", "order_id", order.ID)
	return nil
}
