package main

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/sesaquecruz/goexpert-clean-architecture-lab/config"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/entity"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/event"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/event/handler"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/infra/database"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/infra/grpc/pb"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/infra/grpc/service"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/usecase"
	ev "github.com/sesaquecruz/goexpert-clean-architecture-lab/pkg/event"

	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Configurations
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// Database
	db, err := sql.Open(
		cfg.DBDriver,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName),
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// RabbitMQ
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.AMQPUser, cfg.AMQPPassword, cfg.AMQPHost, cfg.AMQPPort))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	// Repository
	orderRepository := database.NewOrderRepository(db)

	// Event Dispatcher
	eventDispatcher := ev.NewEventDispatcher()

	handler := handler.NewOrderCreatedHandler(ch)
	eventDispatcher.Register(event.EventOrderCreated, handler)

	// UseCase
	createOrder := usecase.NewCreateOrderUseCase(
		&entity.OrderFactory{},
		orderRepository,
		event.EventFactory{},
		eventDispatcher,
	)

	listOrders := usecase.NewListOrdersUseCase(orderRepository)

	// gRCP
	service := service.NewOrderService(createOrder, listOrders)

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, service)
	reflection.Register(grpcServer)

	fmt.Println("gRPC server is running on port", cfg.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCServerPort))
	if err != nil {
		panic(err)
	}

	grpcServer.Serve(lis)
}
