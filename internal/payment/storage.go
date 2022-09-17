package payment

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, p *Payment, id string) (bool, error)
	FindAll(ctx context.Context) ([]Payment, error)
	FindOne(ctx context.Context, id string) (Payment, error)
}
