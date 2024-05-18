package teeworldsecon

import (
	"fmt"
	"regexp"
)

// Function prototype used when an event is handled
type EconEventFunc func(eventPayload string) any

// Represents an event
type EconEvent struct {
	// Event name
	Name string
	// Event regex
	Regex string
	// Event function called when the event is handled
	Func EconEventFunc
	// Event queue fill with the `EconEventFunc` return value when its not nil
	Queue []any
}

// Create a new EconEvent
func NewEconEvent(name string, regex string, f EconEventFunc) *EconEvent {
	return &EconEvent{
		Name:  name,
		Regex: regex,
		Func:  f,
		Queue: []any{},
	}
}

// Event manager controller
type EconEventManager struct {
	Events map[string]*EconEvent
}

// Create an event manager
func NewEconEventManager() *EconEventManager {
	return &EconEventManager{
		Events: make(map[string]*EconEvent),
	}
}

// Returns a `EconEvent` from the `Events` map
func (em *EconEventManager) Event(name string) *EconEvent {
	value, exist := em.Events[name]

	if !exist {
		return nil
	}

	return value
}

// Register an event
func (em *EconEventManager) Register(event *EconEvent) error {
	if event == nil {
		return fmt.Errorf("nil event")
	}

	em.Events[event.Name] = event

	return nil
}

// Call every matching event
func (em *EconEventManager) Call(eventPayload string) {
	for _, event := range em.Events {
		ok, _ := regexp.MatchString(event.Regex, eventPayload)
		if !ok {
			continue
		}

		elem := event.Func(eventPayload)
		if elem != nil {
			event.Queue = append(event.Queue, elem)
		}
	}
}
