package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort           string
	DBUser            string
	DBPassword        string
	DBHost            string
	DBPort            string
	DBName            string
	RabbitMQURL       string
	MidtransServerKey string
	IsProduction      bool
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from system environment variables")
	}

	serverKey := getEnv("MIDTRANS_SERVER_KEY", "")
	
	// Read explicit environment variable to determine mode
	isProdStr := getEnv("MIDTRANS_IS_PRODUCTION", "false")
	isProduction := false
	if strings.ToLower(isProdStr) == "true" {
		isProduction = true
	}

	return &Config{
		AppPort:           getEnv("APP_PORT", "8084"),
		DBUser:            getEnv("DB_USER", "root"),
		DBPassword:        getEnv("DB_PASSWORD", "secretpassword"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBName:            getEnv("DB_NAME", "payment_db"),
		RabbitMQURL:       getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		MidtransServerKey: serverKey,
		IsProduction:      isProduction,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
