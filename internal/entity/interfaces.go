package entity

import "context"

type OrderRepositoryInterface interface {
	Save(ctx context.Context, order Order) error
	FindAll(ctx context.Context) ([]Order, error)
}
