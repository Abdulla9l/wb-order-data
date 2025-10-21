package kafka

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"wb-order-data/internal/models"
	"wb-order-data/internal/service"
)

type Consumer struct {
	service *service.OrderService
}

func NewConsumer(service *service.OrderService) *Consumer {
	return &Consumer{service: service}
}

// Симуляция получения сообщений из Kafka
func (c *Consumer) Start() {
	for {
		// Для демонстрации — случайно получаем один из уже сгенерированных заказов
		order := models.GenerateTestOrder(rand.Intn(10))
		data, _ := json.Marshal(order)

		var o models.Order
		if err := json.Unmarshal(data, &o); err != nil {
			log.Printf("Consumer failed to parse order: %v", err)
			continue
		}

		if err := c.service.SaveOrder(o); err != nil {
			log.Printf("Consumer failed to save order: %v", err)
		} else {
			log.Printf("Consumer saved order: %s", o.OrderUID)
		}

		time.Sleep(5 * time.Second)
	}
}
