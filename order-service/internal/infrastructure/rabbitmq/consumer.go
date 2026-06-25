package rabbitmq

import (
	"context"
	"encoding/json"
	"log/slog"

	"order-service/internal/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventConsumer struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	usecase domain.OrderUsecase
}

func NewRabbitMQConsumer(url string, usecase domain.OrderUsecase) (*EventConsumer, error) {
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
		"order_service_queue", // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,                 // queue name
		"payment.*",            // routing key (listens to all payment events)
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
		"order_service_queue", // queue
		"",                    // consumer
		false,                 // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Info("Shutting down RabbitMQ consumer in Order Service")
				c.ch.Close()
				c.conn.Close()
				return
			case d := <-msgs:
				var event map[string]interface{}
				if err := json.Unmarshal(d.Body, &event); err != nil {
					slog.Error("Failed to unmarshal event", "error", err)
					d.Nack(false, false)
					continue
				}

				if event["event_type"] == "PAYMENT_SUCCESS" {
					payload := event["payload"].(map[string]interface{})
					orderID := payload["order_id"].(string)

					slog.Info("Consumed PAYMENT_SUCCESS event", "order_id", orderID)

					err := c.usecase.MarkOrderAsPaid(context.Background(), orderID)
					if err != nil {
						slog.Error("Failed to mark order as paid", "error", err)
						d.Nack(false, true)
						continue
					}
				} else if event["event_type"] == "PAYMENT_FAILED" {
					payload := event["payload"].(map[string]interface{})
					orderID := payload["order_id"].(string)

					slog.Info("Consumed PAYMENT_FAILED event", "order_id", orderID)

					err := c.usecase.MarkOrderAsFailed(context.Background(), orderID)
					if err != nil {
						slog.Error("Failed to mark order as failed", "error", err)
						d.Nack(false, true)
						continue
					}
				}

				d.Ack(false)
			}
		}
	}()

	slog.Info("Order Service Consumer started, waiting for messages...")
	return nil
}
