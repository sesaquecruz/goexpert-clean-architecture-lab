package event

const (
	EventOrderCreated = "OrderCreated"
)

type Event struct {
	Name    string
	Payload interface{}
}

func (e *Event) GetName() string {
	return e.Name
}

func (e *Event) GetPayload() interface{} {
	return e.Payload
}
