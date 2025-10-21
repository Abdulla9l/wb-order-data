package main

import (
    "log"
    "wb-order-data/internal/config"
    "wb-order-data/internal/database"
    "wb-order-data/internal/cache"
    "wb-order-data/internal/service"
    "wb-order-data/internal/kafka"
    "wb-order-data/internal/server"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal(err)
    }

    db, err := database.NewPostgres(cfg)
    if err != nil {
        log.Fatal(err)
    }

    c := cache.NewCache()

    orderService := service.NewOrderService(db, c)

    kafka.StartConsumer(cfg, orderService)

    server.Start(cfg, orderService)
}
