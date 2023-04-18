package usecase

import (
	"context"

	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/entity"
)

type ListOrdersOutputDTO struct {
	Id         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(orderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: orderRepository,
	}
}

func (u *ListOrdersUseCase) Execute(ctx context.Context) ([]*ListOrdersOutputDTO, error) {
	orders, err := u.OrderRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	output := make([]*ListOrdersOutputDTO, 0)
	for _, order := range orders {
		out := ListOrdersOutputDTO{
			Id:         order.Id.String(),
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}

		output = append(output, &out)
	}

	return output, nil
}
