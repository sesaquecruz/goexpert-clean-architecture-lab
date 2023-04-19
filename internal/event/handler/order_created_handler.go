package handler

import (
	"context"
	"encoding/json"

	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/infra/rabbitmq"
	ev "github.com/sesaquecruz/goexpert-clean-architecture-lab/pkg/event"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderCreatedHandler struct {
	RabbitmqChannel rabbitmq.RabbitMqChannel
}

func NewOrderCreatedHandler(rabbitmqChannel rabbitmq.RabbitMqChannel) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		RabbitmqChannel: rabbitmqChannel,
	}
}

func (h *OrderCreatedHandler) Handle(ctx context.Context, event ev.EventInterface) error {
	payload, err := json.Marshal(event.GetPayload())
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        payload,
	}

	err = h.RabbitmqChannel.PublishWithContext(
		ctx,
		"amq.direct",
		"",
		false,
		false,
		msg,
	)

	return err
}
