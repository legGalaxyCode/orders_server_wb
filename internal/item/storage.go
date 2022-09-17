package item

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, item *Item, orderId string) (bool, error)
	FindAll(ctx context.Context) ([]Item, error)
	FindOne(ctx context.Context, orderId string) (Item, error)
}
