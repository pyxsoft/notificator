package event

import (
	"fmt"
	"reflect"
)

type EventRegistry struct {
	events map[string]reflect.Type
}

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		events: make(map[string]reflect.Type),
	}
}

func (er *EventRegistry) AddEvent(name string, dataType interface{}) {
	er.events[name] = reflect.TypeOf(dataType)
}

func (er *EventRegistry) ValidateEventData(name string, data interface{}) error {
	expectedType, ok := er.events[name]
	if !ok {
		return fmt.Errorf("event type '%s' not registered", name)
	}
	if reflect.TypeOf(data) != expectedType {
		return fmt.Errorf("event data type mismatch for '%s': expected %v, got %v", name, expectedType, reflect.TypeOf(data))
	}
	return nil
}
