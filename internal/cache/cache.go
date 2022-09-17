package cache

import (
	"fmt"
	"sync"
	"test_db_server/internal/order"
)

type Cache struct {
	sync.RWMutex
	Data map[string]order.Order
}

func New() *Cache {
	c := &Cache{
		RWMutex: sync.RWMutex{},
		Data:    make(map[string]order.Order),
	}
	return c
}

func (c *Cache) AddMany(orders ...*order.Order) {
	c.Lock()
	defer c.Unlock()
	for _, ord := range orders {
		c.Data[ord.OrderUid] = *ord
	}
}

func (c *Cache) AddOne(order *order.Order) {
	c.Lock()
	defer c.Unlock()
	c.Data[order.OrderUid] = *order
}

func (c *Cache) GetAll() []order.Order {
	c.RLock()
	defer c.RUnlock()
	orders := make([]order.Order, 0)
	for _, ord := range c.Data {
		orders = append(orders, ord)
	}
	return orders
}

func (c *Cache) GetOne(orderId string) (order.Order, error) {
	c.RLock()
	defer c.RUnlock()
	if ord, ok := c.Data[orderId]; ok {
		return ord, nil
	}
	return order.Order{}, fmt.Errorf("no order with this id: %s", orderId)
}
