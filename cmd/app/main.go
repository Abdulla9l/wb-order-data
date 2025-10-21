package main

import (
	"log"
	"net/http"
	"time"
	"wb-order-data/internal/cache"
	"wb-order-data/internal/kafka"
	"wb-order-data/internal/models"
	"wb-order-data/internal/service"
)

func main() {
	// Создаём кэш
	orderCache := cache.NewOrderCache()

	// Создаём сервис (пока без DB)
	orderService := service.NewOrderService(orderCache)

	// Запускаем продюсера (генерация тестовых заказов)
	producer := kafka.NewProducer(orderService)
	go producer.Start()

	// Запускаем consumer (имитация получения заказов)
	consumer := kafka.NewConsumer(orderService)
	go consumer.Start()

	// HTTP хэндлер для получения заказа по ID
	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}

		order, ok := orderCache.Get(id)
		if !ok {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		data, err := models.MarshalOrder(order)
		if err != nil {
			http.Error(w, "failed to marshal order", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	port := ":8081"
	log.Printf("HTTP server started on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
