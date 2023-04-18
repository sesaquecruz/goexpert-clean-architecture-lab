package usecase

import (
	"context"

	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/entity"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/event"

	ev "github.com/sesaquecruz/goexpert-clean-architecture-lab/pkg/event"
)

type CreateOrderInputDTO struct {
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type CreateOrderOutputDTO struct {
	Id         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderFactory    entity.OrderFactoryInterface
	OrderRepository entity.OrderRepositoryInterface
	EventFactory    event.EventFactory
	EventDispatcher ev.EventDispatcherInterface
}

func NewCreateOrderUseCase(
	orderFactory entity.OrderFactoryInterface,
	orderRepository entity.OrderRepositoryInterface,
	eventFactory event.EventFactory,
	eventDispatcher ev.EventDispatcherInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderFactory:    orderFactory,
		OrderRepository: orderRepository,
		EventFactory:    eventFactory,
		EventDispatcher: eventDispatcher,
	}
}

func (u *CreateOrderUseCase) Execute(ctx context.Context, input CreateOrderInputDTO) (*CreateOrderOutputDTO, error) {
	order, err := u.OrderFactory.NewOrder(input.Price, input.Tax)
	if err != nil {
		return nil, err
	}

	if err := order.CalculateFinalPrice(); err != nil {
		return nil, err
	}

	if err := u.OrderRepository.Save(ctx, order); err != nil {
		return nil, err
	}

	output := CreateOrderOutputDTO{
		Id:         order.Id.String(),
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}

	event := u.EventFactory.NewEvent(event.EventOrderCreated, output)
	u.EventDispatcher.Dispatch(ctx, event)

	return &output, nil
}
