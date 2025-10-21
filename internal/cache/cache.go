package cache

import "wb-order-data/internal/models"

const CacheLimit = 10

type OrderCache struct {
	data   map[string]models.Order
	order  []string
	limit  int
}

func NewOrderCache(limit int) *OrderCache {
	return &OrderCache{
		data:  make(map[string]models.Order),
		order: []string{},
		limit: limit,
	}
}

func (c *OrderCache) Set(orderUID string, order models.Order) {
	c.data[orderUID] = order
	c.order = append(c.order, orderUID)

	if len(c.order) > c.limit {
		oldest := c.order[0]
		c.order = c.order[1:]
		delete(c.data, oldest)
	}
}

func (c *OrderCache) Get(orderUID string) (models.Order, bool) {
	order, ok := c.data[orderUID]
	return order, ok
}
