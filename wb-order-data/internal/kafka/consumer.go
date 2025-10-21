package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"wb-order-data/internal/models"
	"wb-order-data/internal/service"
)

type Consumer struct {
	reader  *kafka.Reader
	service *service.OrderService
}

// NewConsumer создаёт нового Kafka consumer
func NewConsumer(brokers []string, topic, groupID string, orderService *service.OrderService) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		StartOffset: kafka.FirstOffset,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
	})

	return &Consumer{
		reader:  r,
		service: orderService,
	}
}

// Start запускает бесконечное чтение сообщений из Kafka
func (c *Consumer) Start(ctx context.Context) {
	log.Println("Kafka consumer started and waiting for messages...")

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Printf("⚠️  Invalid message format: %v", err)
			continue
		}

		if order.OrderUID == "" {
			log.Printf("⚠️  Skipping invalid order: missing order_uid")
			continue
		}

		if err := c.service.SaveOrder(ctx, order); err != nil {
			log.Printf("❌ Failed to save order %s: %v", order.OrderUID, err)
			continue
		}

		log.Printf("✅ Order %s processed and saved", order.OrderUID)
		time.Sleep(200 * time.Millisecond) // лёгкая пауза, чтобы не спамить лог
	}
}

// Close закрывает подключение к Kafka
func (c *Consumer) Close() {
	if err := c.reader.Close(); err != nil {
		log.Printf("Error closing Kafka reader: %v", err)
	}
}
