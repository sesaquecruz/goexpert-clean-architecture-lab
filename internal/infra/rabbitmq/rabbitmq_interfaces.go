package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqChannel interface {
	PublishWithContext(
		ctx context.Context,
		exchange string,
		key string,
		mandatory bool,
		immediate bool,
		msg amqp.Publishing,
	) error

	Consume(queue string,
		consumer string,
		autoAck bool,
		exclusive bool,
		noLocal bool,
		noWait bool,
		args amqp.Table,
	) (<-chan amqp.Delivery, error)
}
