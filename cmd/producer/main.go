package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

func main() {
	godotenv.Load()

	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	log.Println("Starting Kafka producer - generating random orders every 10 seconds...")

	for {
		order := generateRandomOrder()

		jsonData, err := json.Marshal(order)
		if err != nil {
			log.Printf("Error marshaling order: %v", err)
			continue
		}

		orderUID := order["order_uid"].(string)
		message := kafka.Message{
			Key:   []byte(orderUID),
			Value: jsonData,
		}

		err = writer.WriteMessages(context.Background(), message)
		if err != nil {
			log.Printf("Error sending message to Kafka: %v", err)
		} else {
			log.Printf(" New order generated: %s", orderUID)
		}

		time.Sleep(10 * time.Second)
	}
}

func generateRandomOrder() map[string]interface{} {
	timestamp := time.Now().Unix()
	randomNum := rand.Intn(10000)
	orderUID := fmt.Sprintf("order-%d-%d", timestamp, randomNum)

	return map[string]interface{}{
		"order_uid":    orderUID,
		"track_number": fmt.Sprintf("TRACK-%d", randomNum),
		"entry":        fmt.Sprintf("ENTRY-%s", randomString(4)),
		"delivery": map[string]string{
			"name":    fmt.Sprintf("user_%s", randomString(6)),
			"phone":   fmt.Sprintf("+7%s", randomDigits(10)),
			"zip":     randomDigits(6),
			"city":    fmt.Sprintf("city_%s", randomString(5)),
			"address": fmt.Sprintf("address_%s_%d", randomString(4), randomNum),
			"region":  fmt.Sprintf("region_%s", randomString(4)),
			"email":   fmt.Sprintf("email_%s@test.com", randomString(6)),
		},
		"payment": map[string]interface{}{
			"transaction":   fmt.Sprintf("trans_%s", randomString(8)),
			"request_id":    fmt.Sprintf("req_%s", randomString(6)),
			"currency":      "USD",
			"provider":      fmt.Sprintf("provider_%s", randomString(4)),
			"amount":        rand.Intn(10000) + 1000,
			"payment_dt":    timestamp,
			"bank":          fmt.Sprintf("bank_%s", randomString(4)),
			"delivery_cost": rand.Intn(2000) + 500,
			"goods_total":   rand.Intn(8000) + 1000,
			"custom_fee":    rand.Intn(100),
		},
		"items": []map[string]interface{}{
			{
				"chrt_id":      rand.Intn(10000000),
				"track_number": fmt.Sprintf("ITEM-TRACK-%d", randomNum),
				"price":        rand.Intn(5000) + 500,
				"rid":          fmt.Sprintf("rid_%s", randomString(10)),
				"name":         fmt.Sprintf("product_%s", randomString(6)),
				"sale":         rand.Intn(50),
				"size":         fmt.Sprintf("SIZE-%s", randomString(2)),
				"total_price":  rand.Intn(5000) + 500,
				"nm_id":        rand.Intn(10000000),
				"brand":        fmt.Sprintf("brand_%s", randomString(5)),
				"status":       rand.Intn(400) + 200,
			},
		},
		"locale":             "en",
		"internal_signature": fmt.Sprintf("sig_%s", randomString(8)),
		"customer_id":        fmt.Sprintf("customer_%s", randomString(6)),
		"delivery_service":   fmt.Sprintf("service_%s", randomString(4)),
		"shardkey":           fmt.Sprintf("%d", rand.Intn(10)),
		"sm_id":              rand.Intn(100),
		"date_created":       time.Now().Format(time.RFC3339),
		"oof_shard":          fmt.Sprintf("%d", rand.Intn(3)+1),
	}
}

func randomString(length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func randomDigits(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}
