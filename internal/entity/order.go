package entity

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrOrderInvalidId    = errors.New("invalid id")
	ErrOrderInvalidPrice = errors.New("invalid price")
	ErrOrderInvalidTax   = errors.New("invalid tax")
)

type Order struct {
	Id         uuid.UUID
	Price      float64
	Tax        float64
	FinalPrice float64
}

func (o *Order) CalculateFinalPrice() error {
	o.FinalPrice = o.Price + o.Tax
	return o.IsValid()
}

func (o *Order) IsValid() error {
	if o.Id == uuid.Nil {
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
