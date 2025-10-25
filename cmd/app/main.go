package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wb-order-data/internal/cache"
	"wb-order-data/internal/kafka"
	"wb-order-data/internal/repository"
	"wb-order-data/internal/server"
	"wb-order-data/internal/service"
)

func main() {
	kafkaBrokers := []string{"localhost:9092"}
	kafkaTopic := "orders"
	consumerGroupID := "order-service-group"

	db, err := repository.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	orderRepo := repository.NewOrderRepository(db)

	orderCache := cache.NewOrderCache(10)

	orderService := service.NewOrderService(orderRepo, orderCache)

	log.Println("Restoring cache from database...")
	if err := orderService.RestoreCacheFromDB(); err != nil {
		log.Printf("Warning: failed to restore cache: %v", err)
	} else {
		log.Println("Cache successfully restored from database")
	}

	kafkaConsumer := kafka.NewConsumer(kafkaBrokers, kafkaTopic, consumerGroupID, orderService)
	defer kafkaConsumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		kafkaConsumer.Start(ctx)
	}()

	httpServer := server.NewHTTPServer(orderService)

	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := httpServer.Start(":8080"); err != nil {
			log.Fatal("Failed to start HTTP server:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")

	cancel()
	time.Sleep(2 * time.Second)

	log.Println("Server stopped")
}
