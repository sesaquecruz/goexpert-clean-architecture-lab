package event

type EventFactoryInterface interface {
	NewEvent(name string, payload interface{}) *Event
}
