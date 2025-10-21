package kafka

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"wb-order-data/internal/models"
	"wb-order-data/internal/service"
)

type Producer struct {
	service *service.OrderService
}

func NewProducer(service *service.OrderService) *Producer {
	return &Producer{service: service}
}

// Имитируем генерацию новых заказов каждые 5 секунд
func (p *Producer) Start() {
	id := 0
	for {
		id++
		order := models.GenerateTestOrder(id)

		data, err := json.Marshal(order)
		if err != nil {
			log.Printf("Failed to marshal order: %v", err)
			continue
		}

		// Отправляем заказ в consumer (в реальности — Kafka topic)
		log.Printf("Producer sent order: %s", order.OrderUID)
		go func(msg []byte) {
			// Для теста можно напрямую вызвать consumer через сервис
			var o models.Order
			if err := json.Unmarshal(msg, &o); err != nil {
				log.Printf("Consumer parse error: %v", err)
				return
			}
			if err := p.service.SaveOrder(o); err != nil {
				log.Printf("Consumer save error: %v", err)
			}
		}(data)

		time.Sleep(5 * time.Second)
	}
}
