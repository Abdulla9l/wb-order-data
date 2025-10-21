package cache

import "wb-order-data/internal/model"

const CacheLimit = 10

var CacheMap = make(map[string]model.Order)
var CacheOrder []string

func AddOrder(order model.Order) {
    CacheMap[order.OrderUID] = order
    CacheOrder = append(CacheOrder, order.OrderUID)
    if len(CacheOrder) > CacheLimit {
        oldest := CacheOrder[0]
        CacheOrder = CacheOrder[1:]
        delete(CacheMap, oldest)
    }
}

func GetOrder(id string) (model.Order, bool) {
    order, ok := CacheMap[id]
    return order, ok
}
