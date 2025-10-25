package service

import (
	"testing"
	"wb-order-data/internal/cache"
	"wb-order-data/internal/repository"
)

func TestOrderService_GetOrder(t *testing.T) {
	// Используем реальный репозиторий с nil DB для простых тестов
	repo := repository.NewOrderRepository(nil)
	cache := cache.NewOrderCache(10)
	service := NewOrderService(repo, cache)

	order, err := service.GetOrder("nonexistent")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if order != nil {
		t.Error("Expected nil for nonexistent order")
	}
}

func TestOrderService_RestoreCache(t *testing.T) {
	repo := repository.NewOrderRepository(nil)
	cache := cache.NewOrderCache(10)
	service := NewOrderService(repo, cache)

	err := service.RestoreCacheFromDB()
	if err != nil {
		t.Logf("Expected error with nil DB: %v", err)
	}
}
