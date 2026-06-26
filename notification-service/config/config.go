package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQURL string
	Port        string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load() // ignore error, might be running in docker

	rmq := os.Getenv("RABBITMQ_URL")
	if rmq == "" {
		rmq = "amqp://guest:guest@localhost:5672/"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	return &Config{
		RabbitMQURL: rmq,
		Port:        port,
	}, nil
}
