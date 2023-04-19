package event

import "context"

type EventInterface interface {
	GetName() string
	GetPayload() interface{}
}

type EventHandlerInterface interface {
	Handle(ctx context.Context, event EventInterface) error
}

type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error
	Remove(eventName string, handler EventHandlerInterface) error
	Has(eventName string, handler EventHandlerInterface) bool
	Dispatch(ctx context.Context, event EventInterface) []error
	Clear()
}
