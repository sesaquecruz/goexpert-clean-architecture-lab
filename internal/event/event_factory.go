package event

type EventFactory struct{}

func (f *EventFactory) NewEvent(name string, payload interface{}) *Event {
	return &Event{
		Name:    name,
		Payload: payload,
	}
}
