package config

import (
	"os"
)

type Config struct {
	KafkaBrokers  []string
	KafkaTopic    string
	ConsumerGroup string
	DatabaseURL   string
	HTTPPort      string
}

func Load() *Config {
	return &Config{
		KafkaBrokers:  getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}, ","),
		KafkaTopic:    getEnv("KAFKA_TOPIC", "orders"),
		ConsumerGroup: getEnv("KAFKA_CONSUMER_GROUP", "order-service"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/ordersdb"),
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string, separator string) []string {
	if value := os.Getenv(key); value != "" {
		var result []string
		return result
	}
	return defaultValue
}
