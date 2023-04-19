package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/sesaquecruz/goexpert-clean-architecture-lab/config"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/entity"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/event"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/event/handler"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/infra/database"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/infra/graphql"
	gql_resolver "github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/infra/graphql/resolver"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/infra/grpc/pb"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/infra/grpc/service"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/infra/web"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/usecase"
	ev "github.com/sesaquecruz/goexpert-clean-architecture-lab/pkg/event"

	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/go-chi/chi/v5"

	gql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	// Configurations
	cfg, err := config.LoadConfig()
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

	err = waitForMySQL(context.Background(), db)
	if err != nil {
		panic(err)
	}

	// RabbitMQ
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.AMQPUser, cfg.AMQPPassword, cfg.AMQPHost, cfg.AMQPPort))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	err = waitForRabbitMQ(context.Background(), conn)
	if err != nil {
		panic(err)
	}

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

	// UseCases
	createOrderUseCase := usecase.NewCreateOrderUseCase(
		&entity.OrderFactory{},
		orderRepository,
		event.EventFactory{},
		eventDispatcher,
	)

	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)

	//
	// Services
	//
	wg := &sync.WaitGroup{}

	// gRCP
	service := service.NewOrderService(createOrderUseCase, listOrdersUseCase)

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, service)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCServerPort))
	if err != nil {
		panic(err)
	}

	fmt.Println("gRPC server is running on port", cfg.GRPCServerPort)
	wg.Add(1)
	go func() {
		defer wg.Done()
		grpcServer.Serve(lis)
	}()

	// Rest
	orderWebHandlers := web.NewOrderWebHandlers(createOrderUseCase, listOrdersUseCase)

	router := chi.NewRouter()
	router.Route("/order", func(r chi.Router) {
		r.Post("/", orderWebHandlers.CreateOrderHandler)
		r.Get("/", orderWebHandlers.ListOrdersHandler)
	})

	fmt.Println("Rest server is running on port", cfg.RESTServerPort)
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.ListenAndServe(fmt.Sprintf(":%s", cfg.RESTServerPort), router)
	}()

	// GraphQL
	graphQLServer := gql_handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &gql_resolver.Resolver{
			CreateOrderUseCase: createOrderUseCase,
			ListOrdersUseCase:  listOrdersUseCase,
		},
	}))

	graphQLHttpServer := http.NewServeMux()
	graphQLHttpServer.Handle("/", playground.Handler("GraphQL Playground", "/order"))
	graphQLHttpServer.Handle("/order", graphQLServer)

	fmt.Println("GraphQL server is running on port", cfg.GRAPHQLServerPort)
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.ListenAndServe(fmt.Sprintf(":%s", cfg.GRAPHQLServerPort), graphQLHttpServer)
	}()

	// Waiting services
	wg.Wait()
}

func waitForMySQL(ctx context.Context, db *sql.DB) error {
	for i := 0; i <= 30; {
		if err := db.PingContext(ctx); err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
			i++
		}
	}

	return errors.New("unable to connect to the MySQL")
}

func waitForRabbitMQ(ctx context.Context, conn *amqp.Connection) error {
	for i := 0; i <= 30; {
		if ch, err := conn.Channel(); err == nil {
			ch.Close()
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
			i++
		}
	}

	return errors.New("unable to connect to the RabbitMQ")
}
