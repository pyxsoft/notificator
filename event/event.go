package event

type Event struct {
	Name string
	Data interface{}
}

func NewEvent(name string, data interface{}) *Event {
	return &Event{
		Name: name,
		Data: data,
	}
}

func (e *Event) GetData() interface{} {
	return e.Data
}

func (e *Event) GetType() string {
	return e.Name
}
