package service

import (
	"context"
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

func (s *OrderService) SaveOrder(ctx context.Context, order models.Order) error {
	if err := s.repo.SaveOrder(ctx, order); err != nil {
		log.Printf("❌ Error saving order to DB: %v", err)
		return err
	}

	s.cache.Set(order.OrderUID, order)

	log.Printf("✅ Order %s saved to DB and cache", order.OrderUID)
	return nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (models.Order, bool, error) {

	if order, found := s.cache.Get(id); found {
		return order, true, nil
	}

	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return models.Order{}, false, err
	}

	s.cache.Set(id, order)

	return order, false, nil
}
