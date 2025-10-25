package cache

import (
	"testing"
	"wb-order-data/internal/models"
)

func TestOrderCache_SetAndGet(t *testing.T) {
	cache := NewOrderCache(10)
	order := models.Order{OrderUID: "test-1", TrackNumber: "track-1"}
	cache.Set(order)

	retrieved, ok := cache.Get("test-1")
	if !ok {
		t.Error("Failed to get cached order")
	}
	if retrieved.OrderUID != "test-1" {
		t.Errorf("Expected order ID 'test-1', got '%s'", retrieved.OrderUID)
	}

	_, ok = cache.Get("nonexistent")
	if ok {
		t.Error("Should not find non-existent order")
	}
}

func TestOrderCache_LRU_Eviction(t *testing.T) {
	cache := NewOrderCache(2)

	cache.Set(models.Order{OrderUID: "order-1"})
	cache.Set(models.Order{OrderUID: "order-2"})
	cache.Set(models.Order{OrderUID: "order-3"})

	_, ok := cache.Get("order-1")
	if ok {
		t.Error("order-1 should be evicted from cache")
	}

	_, ok = cache.Get("order-2")
	if !ok {
		t.Error("order-2 should be in cache")
	}

	_, ok = cache.Get("order-3")
	if !ok {
		t.Error("order-3 should be in cache")
	}
}

func TestOrderCache_Size(t *testing.T) {
	cache := NewOrderCache(5)

	if cache.Size() != 0 {
		t.Error("New cache should be empty")
	}

	cache.Set(models.Order{OrderUID: "order-1"})
	cache.Set(models.Order{OrderUID: "order-2"})

	if cache.Size() != 2 {
		t.Errorf("Expected cache size 2, got %d", cache.Size())
	}
}
