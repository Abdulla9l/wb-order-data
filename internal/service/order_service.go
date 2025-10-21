package service

import (
	"wb-order-data/internal/cache"
	"wb-order-data/internal/models"
	"wb-order-data/internal/repository"
)

type OrderService struct {
	repo  *repository.OrderRepository
	cache *cache.OrderCache
}

func NewOrderService(repo *repository.OrderRepository, c *cache.OrderCache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: c,
	}
}

func (s *OrderService) SaveOrder(order models.Order) error {
	// Сохраняем в БД через указатель
	if err := s.repo.SaveOrder(&order); err != nil {
		return err
	}

	// Сохраняем в кэше по значению
	s.cache.Set(order.OrderUID, order)
	return nil
}

func (s *OrderService) GetOrder(id string) (models.Order, bool, error) {
	// Сначала ищем в кэше
	if order, found := s.cache.Get(id); found {
		return order, true, nil
	}

	// Извлекаем из БД
	orderPtr, err := s.repo.GetOrderByID(id) // возвращает *models.Order
	if err != nil {
		return models.Order{}, false, err
	}

	// Сохраняем в кэше по значению
	s.cache.Set(orderPtr.OrderUID, *orderPtr)

	return *orderPtr, false, nil
}
