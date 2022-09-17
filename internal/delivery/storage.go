package delivery

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, del *Delivery, orderId string) (bool, error)
	FindAll(ctx context.Context) ([]Delivery, error)
	FindOne(ctx context.Context, id string) (Delivery, error)
}
