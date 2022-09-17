package order

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, order *Order) (bool, error)
	CreateMany(ctx context.Context, orders ...*Order) (bool, error)
	FindAll(ctx context.Context) ([]Order, error)
	FindOne(ctx context.Context, orderId string) (Order, error)
}
