# Clean Architecture Lab - Go Expert

This project is designed to streamline the management of purchase orders. It follows a clear architecture and utilizes MySQL as the database. It supports creating new orders and listing existing orders through gRPC, REST API or GraphQL interfaces. Additionally, it integrates with RabbitMQ to send order data to a designated exchange upon the creation of a new order.

## Requirements

To use this program, you will need:

- Docker
- A stable internet connection

## Installation

1. Clone this repository:

```
git clone https://github.com/sesaquecruz/goexpert-clean-architecture-lab
```

2. Enter the project directory:

```
cd goexpert-clean-architecture-lab
```

3. Run the docker compose:

```
docker compose up --build
```

## Usage

### RabbitMQ

1. Access and use (user, user) as credentials:

```
http://localhost:15672
```

2. Create a new queue with any name. E.g. **orders**.

3. Bind the previously created queue with the exchange **amq.direct** using **order** as routing key.

### gRPC

1. Use a gRPC client to connect on port 50051. Using [evans](https://github.com/ktr0731/evans) for example, run:

```
evans -r repl
```

2. To list all created orders, run:

```
call ListOrders
```

3. To create a new order, run and enter the asked values:

```
call CreateOrder
```

### REST API

1. To list all created orders, run:

```
curl -v http://localhost:8080/order
```

2. To create a new order, run:

```
curl -v -X POST http://localhost:8080/order -H "Content-Type: application/json" -d '{"Price": 23.00, "Tax": 0.23}'
```

### GraphQL

1. Access the GraphQL Playground:

```
http://localhost:8000
```

1. To list all created orders, run:

```
query listOrders {
  listOrders {
    Id
    Price
    Tax
    FinalPrice
  }
}
```

2. To create a new order, run:

```
mutation createOrder {
  createOrder(input: {Price: 42.00, Tax: 0.42})
  {
    Id
    Price
    Tax
    FinalPrice
  }
}
```

## Troubleshooting

See [docker-compose.yml](./docker-compose.yml) to verify or change the services and port values.

## License

This project is licensed under the MIT License. See [LICENSE](./LICENSE) file for more information.
