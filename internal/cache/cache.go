package cache

import (
	"sync"
	"wb-order-data/internal/models"
)

type OrderCache struct {
	data  map[string]models.Order
	order []string
	mutex sync.RWMutex
	limit int
}

func NewOrderCache(limit int) *OrderCache {
	return &OrderCache{
		data:  make(map[string]models.Order),
		order: make([]string, 0),
		limit: limit,
	}
}

func (c *OrderCache) Set(order models.Order) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i, id := range c.order {
		if id == order.OrderUID {
			c.order = append(c.order[:i], c.order[i+1:]...)
			break
		}
	}

	c.data[order.OrderUID] = order
	c.order = append(c.order, order.OrderUID)

	if len(c.order) > c.limit {
		oldest := c.order[0]
		c.order = c.order[1:]
		delete(c.data, oldest)
	}
}

func (c *OrderCache) Get(id string) (models.Order, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	order, ok := c.data[id]
	return order, ok
}

func (c *OrderCache) GetAll() map[string]models.Order {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	result := make(map[string]models.Order)
	for k, v := range c.data {
		result[k] = v
	}
	return result
}

func (c *OrderCache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.data)
}
