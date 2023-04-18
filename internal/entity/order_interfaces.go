package entity

import "context"

type OrderFactoryInterface interface {
	NewOrder(price float64, tax float64) (*Order, error)
}

type OrderRepositoryInterface interface {
	Save(ctx context.Context, order *Order) error
	FindAll(ctx context.Context) ([]*Order, error)
}
