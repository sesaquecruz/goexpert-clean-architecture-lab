package handler

import (
	"context"
	"encoding/json"

	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/infra/rabbitmq"
	ev "github.com/sesaquecruz/goexpert-clean-architecture-lab/pkg/event"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderCreatedHandler struct {
	RabbitmqChannel rabbitmq.RabbitMqChannelPublishInterface
	Exchange        string
	Key             string
}

func NewOrderCreatedHandler(
	rabbitmqChannel rabbitmq.RabbitMqChannelPublishInterface,
	exchange string,
	key string,
) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		RabbitmqChannel: rabbitmqChannel,
		Exchange:        exchange,
		Key:             key,
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
		h.Exchange,
		h.Key,
		false,
		false,
		msg,
	)

	return err
}
