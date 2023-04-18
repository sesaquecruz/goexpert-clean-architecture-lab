package event

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrorHandlerAlreadyRegistered = errors.New("handler already registered")
	ErrorEventNotFound            = errors.New("event not found")
	ErrorHandlerNotFound          = errors.New("handler not found")
)

type EventDispatcher struct {
	Handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		Handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if handlers, ok := ed.Handlers[eventName]; ok {
		for _, h := range handlers {
			if h == handler {
				return ErrorHandlerAlreadyRegistered
			}
		}
	}

	ed.Handlers[eventName] = append(ed.Handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if handlers, ok := ed.Handlers[eventName]; ok {
		for i, h := range handlers {
			if h == handler {
				ed.Handlers[eventName] = append(handlers[:i], handlers[i+1:]...)
				return nil
			}
		}
	} else {
		return ErrorEventNotFound
	}

	return ErrorHandlerNotFound
}

func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if handlers, ok := ed.Handlers[eventName]; ok {
		for _, h := range handlers {
			if h == handler {
				return true
			}
		}
	}

	return false
}

func (ed *EventDispatcher) Dispatch(ctx context.Context, event EventInterface) []error {
	if handlers, ok := ed.Handlers[event.GetName()]; ok {
		ch := make(chan error, len(handlers))
		wg := &sync.WaitGroup{}

		for _, h := range handlers {
			eh := h
			wg.Add(1)

			go func() {
				err := eh.Handle(ctx, event)
				if err != nil {
					ch <- err
				}

				wg.Done()
			}()
		}

		wg.Wait()
		close(ch)

		var errs []error

		for err := range ch {
			errs = append(errs, err)
		}

		if len(errs) > 0 {
			return errs
		}

		return nil
	}

	return []error{ErrorEventNotFound}
}

func (ed *EventDispatcher) Clear() {
	ed.Handlers = make(map[string][]EventHandlerInterface)
}
