package kafka

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "time"

    "github.com/segmentio/kafka-go"
    "wb-order-data/internal/config"
    "wb-order-data/internal/models"
)

func StartProducer(cfg *config.Config) {
    writer := &kafka.Writer{
        Addr:     kafka.TCP(cfg.KafkaURL),
        Topic:    "orders",
        Balancer: &kafka.LeastBytes{},
    }
    defer writer.Close()

    log.Println("Kafka producer started...")

    for {
        order := generateRandomOrder()
        data, _ := json.Marshal(order)

        err := writer.WriteMessages(context.Background(),
            kafka.Message{
                Key:   []byte(order.OrderUID),
                Value: data,
            },
        )
        if err != nil {
            log.Println("Error sending order:", err)
        } else {
            log.Println("Sent order", order.OrderUID)
        }

        time.Sleep(5 * time.Second)
    }
}

func generateRandomOrder() models.Order {
    uid := fmt.Sprintf("test-order-%d", rand.Intn(1000))
    return models.Order{
        OrderUID: uid,
        TrackNumber: fmt.Sprintf("TRACK-%d", rand.Intn(99999)),
        Entry: "WBIL",
        DateCreated: time.Now(),
        // delivery, payment, items можно добавить по аналогии
    }
}
