package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"wb-order-data/internal/models"
	"wb-order-data/internal/service"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	service *service.OrderService
}

func NewConsumer(brokers []string, topic string, groupID string, service *service.OrderService) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	return &Consumer{
		reader:  reader,
		service: service,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Printf("Starting Kafka consumer for topic: %s", c.reader.Config().Topic)

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Kafka consumer...")
			c.reader.Close()
			return
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				log.Printf("Error fetching message: %v", err)
				continue
			}

			log.Printf("Received message from Kafka: key=%s, partition=%d, offset=%d",
				string(msg.Key), msg.Partition, msg.Offset)

			if err := c.processMessage(msg.Value); err != nil {
				log.Printf("Error processing message: %v", err)
				continue
			}

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Error committing message: %v", err)
			} else {
				log.Printf("Successfully processed and committed message: %s", string(msg.Key))
			}
		}
	}
}

func (c *Consumer) processMessage(data []byte) error {
	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return fmt.Errorf("json unmarshal error: %v", err)
	}

	if order.OrderUID == "" {
		return fmt.Errorf("invalid order: missing OrderUID")
	}

	if err := c.service.ProcessOrderFromMessage(data); err != nil {
		return fmt.Errorf("service processing error: %v", err)
	}

	return nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
