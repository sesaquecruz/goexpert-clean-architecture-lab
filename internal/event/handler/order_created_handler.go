package handler

import (
	"context"
	"encoding/json"

	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/event"
	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/infra/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderCreatedHandler struct {
	rabbitmqChannel rabbitmq.RabbitMqChannel
}

func (h *OrderCreatedHandler) Handle(ctx context.Context, event event.Event) error {
	payload, err := json.Marshal(event.GetPayload())
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        payload,
	}

	err = h.rabbitmqChannel.PublishWithContext(
		ctx,
		"amq.direct",
		"",
		false,
		false,
		msg,
	)

	return err
}
