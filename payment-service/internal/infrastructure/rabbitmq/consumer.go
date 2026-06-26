package rabbitmq

import (
	"context"
	"encoding/json"
	"log/slog"

	"payment-service/internal/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventConsumer struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	usecase domain.PaymentUsecase
}

func NewRabbitMQConsumer(url string, usecase domain.PaymentUsecase) (*EventConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

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

	q, err := ch.QueueDeclare(
		"payment_service_queue", // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,                 // queue name
		"order.created",        // routing key
		"food_delivery_events", // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &EventConsumer{
		conn:    conn,
		ch:      ch,
		usecase: usecase,
	}, nil
}

func (c *EventConsumer) StartConsuming(ctx context.Context) error {
	msgs, err := c.ch.Consume(
		"payment_service_queue", // queue
		"",                      // consumer
		false,                   // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Info("Shutting down RabbitMQ consumer")
				c.ch.Close()
				c.conn.Close()
				return
			case d := <-msgs:
				// Parse event
				var event map[string]interface{}
				if err := json.Unmarshal(d.Body, &event); err != nil {
					slog.Error("Failed to unmarshal event", "error", err)
					d.Nack(false, false)
					continue
				}

				if event["event_type"] == "ORDER_CREATED" {
					payload := event["payload"].(map[string]interface{})
					orderID := payload["order_id"].(string)
					customerID := payload["customer_id"].(string)
					customerEmail := ""
					if val, ok := payload["customer_email"].(string); ok {
						customerEmail = val
					}
					amount := payload["total_amount"].(float64)

					slog.Info("Consumed ORDER_CREATED event", "order_id", orderID)

					// Process
					err := c.usecase.ProcessOrderCreated(context.Background(), orderID, customerID, customerEmail, amount)
					if err != nil {
						slog.Error("Failed to process order created", "error", err)
						// Might want to retry by Nacking, but for prototyping we'll just log
						d.Nack(false, true)
						continue
					}
				}

				// Acknowledge the message
				d.Ack(false)
			}
		}
	}()

	slog.Info("RabbitMQ Consumer started, waiting for messages...")
	return nil
}
