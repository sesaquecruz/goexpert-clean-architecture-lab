package entity

import "github.com/google/uuid"

type OrderFactory struct{}

func (f *OrderFactory) NewOrder(price float64, tax float64) (*Order, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	order := &Order{
		Id:    id,
		Price: price,
		Tax:   tax,
	}

	if err := order.IsValid(); err != nil {
		return nil, err
	}

	return order, nil
}
