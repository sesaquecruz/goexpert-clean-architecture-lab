package entity

import "errors"

var (
	ErrOrderInvalidId    = errors.New("invalid id")
	ErrOrderInvalidPrice = errors.New("invalid price")
	ErrOrderInvalidTax   = errors.New("invalid tax")
)

type Order struct {
	ID         string
	Price      float64
	Tax        float64
	FinalPrice float64
}

func NewOrder(id string, price float64, tax float64) (*Order, error) {
	order := &Order{
		ID:    id,
		Price: price,
		Tax:   tax,
	}
	err := order.IsValid()
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *Order) CalculateFinalPrice() error {
	o.FinalPrice = o.Price + o.Tax
	return o.IsValid()
}

func (o *Order) IsValid() error {
	if o.ID == "" {
		return ErrOrderInvalidId
	}
	if o.Price <= 0 {
		return ErrOrderInvalidPrice
	}
	if o.Tax <= 0 {
		return ErrOrderInvalidTax
	}
	return nil
}
