package service

import (
	"encoding/json"
	"fmt"
	"log"
	"wb-order-data/internal/cache"
	"wb-order-data/internal/models"
	"wb-order-data/internal/repository"
)

type OrderService struct {
	repo  *repository.OrderRepository
	cache *cache.OrderCache
}

func NewOrderService(repo *repository.OrderRepository, cache *cache.OrderCache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}

func (s *OrderService) GetOrder(orderUID string) (*models.Order, error) {
	return s.GetOrderByID(orderUID)
}

func (s *OrderService) GetOrderByID(orderUID string) (*models.Order, error) {
	if order, ok := s.cache.Get(orderUID); ok {
		log.Printf("Order %s found in cache", orderUID)
		return &order, nil
	}

	log.Printf("Order %s not in cache, querying database", orderUID)

	order, err := s.repo.GetOrderByID(orderUID)
	if err != nil {
		return nil, fmt.Errorf("db query error: %v", err)
	}

	if order != nil {
		s.cache.Set(*order)
		log.Printf("Order %s loaded from DB and cached", orderUID)
	} else {
		log.Printf("Order %s not found in database", orderUID)
	}

	return order, nil
}

func (s *OrderService) ProcessOrderFromMessage(orderData []byte) error {
	var order models.Order
	if err := json.Unmarshal(orderData, &order); err != nil {
		log.Printf("Error unmarshaling order: %v", err)
		return fmt.Errorf("json unmarshal error: %v", err)
	}

	if order.OrderUID == "" {
		log.Printf("Invalid order: missing OrderUID")
		return fmt.Errorf("invalid order: missing OrderUID")
	}

	if err := s.repo.SaveOrder(&order); err != nil {
		log.Printf("Error saving order to DB: %v", err)
		return fmt.Errorf("db save error: %v", err)
	}

	s.cache.Set(order)
	log.Printf("Order %s processed and cached (Kafka -> DB -> Cache)", order.OrderUID)

	return nil
}

func (s *OrderService) RestoreCacheFromDB() error {
	log.Println("Restoring cache from database...")

	orders, err := s.repo.GetAllOrders()
	if err != nil {
		return fmt.Errorf("failed to get orders from DB: %v", err)
	}

	cachedCount := 0
	for _, order := range orders {
		s.cache.Set(order)
		cachedCount++
	}

	log.Printf("Cache restored with %d orders (cache size: %d)", cachedCount, s.cache.Size())
	return nil
}

// GetCacheStats возвращает статистику кэша
func (s *OrderService) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"cache_size":  s.cache.Size(),
		"cache_limit": 10,
	}
}
