package service

import (
	"context"

	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/infra/grpc/pb"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrderUseCase   *usecase.ListOrdersUseCase
}

func NewOrderService(
	createOrderUseCase *usecase.CreateOrderUseCase,
	listOrderUseCase *usecase.ListOrdersUseCase,
) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrderUseCase:   listOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderInput) (*pb.Order, error) {
	input := usecase.CreateOrderInputDTO{
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}

	output, err := s.CreateOrderUseCase.Execute(ctx, input)
	if err != nil {
		return nil, err
	}

	order := pb.Order{
		Id:         output.Id,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}

	return &order, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.Empty) (*pb.Orders, error) {
	output, err := s.ListOrderUseCase.Execute(ctx)
	if err != nil {
		return nil, err
	}

	var orders pb.Orders
	for _, out := range output {
		order := pb.Order{
			Id:         out.Id,
			Price:      float32(out.Price),
			Tax:        float32(out.Tax),
			FinalPrice: float32(out.FinalPrice),
		}

		orders.Orders = append(orders.Orders, &order)
	}

	return &orders, nil
}
