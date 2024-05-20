package teeworldsecon

import (
	"fmt"
	"regexp"
	"sync"
)

// Represents an entry for the event manager registry
type EconEventManagerEntry struct {
	// Econ event
	event *EconEvent
	// `event.Func` return value is added in `queue` if not equal to nil
	queue []any
	// Mutex to protect the queue
	mu sync.Mutex
}

// Pop an element from the queue
func (ee *EconEventManagerEntry) QueuePop(index int) (any, error) {
	ee.mu.Lock()
	defer ee.mu.Unlock()

	queueLen := len(ee.queue)

	if index < 0 || index >= queueLen {
		return nil, fmt.Errorf("invalid index")
	}

	ret := ee.queue[index]
	ee.queue = append(ee.queue[:index], ee.queue[index+1:]...)

	return ret, nil
}

// Get the queue size
func (ee *EconEventManagerEntry) QueueSize() int {
	ee.mu.Lock()
	defer ee.mu.Unlock()

	return len(ee.queue)
}

// Get an element from the queue
func (ee *EconEventManagerEntry) QueueGet(index int) (any, error) {
	ee.mu.Lock()
	defer ee.mu.Unlock()

	if index < 0 || index >= len(ee.queue) {
		return nil, fmt.Errorf("invalid index")
	}

	return ee.queue[index], nil
}

// Create a new entry in the event manager registry
func NewEconEventManagerEntry(event *EconEvent) *EconEventManagerEntry {
	return &EconEventManagerEntry{
		event: event,
		queue: []any{},
	}
}

// Event manager controller
type EconEventManager struct {
	events map[string]*EconEventManagerEntry
	mu     sync.Mutex
}

// Create an event manager
func NewEconEventManager() *EconEventManager {
	return &EconEventManager{
		events: make(map[string]*EconEventManagerEntry),
	}
}

// Returns an event entry for the events registry
func (em *EconEventManager) Entry(name string) *EconEventManagerEntry {
	em.mu.Lock()
	defer em.mu.Unlock()

	return em.events[name]
}

// Register an event
func (em *EconEventManager) Register(event *EconEvent) error {
	if event == nil {
		return fmt.Errorf("nil event")
	}

	em.mu.Lock()
	defer em.mu.Unlock()

	em.events[event.Name] = NewEconEventManagerEntry(event)

	return nil
}

// Delete an event
func (em *EconEventManager) Delete(name string) {
	em.mu.Lock()
	defer em.mu.Unlock()

	delete(em.events, name)
}

// Call every matching event
func (em *EconEventManager) Handle(eventPayload string) {
	em.mu.Lock()
	defer em.mu.Unlock()

	for _, eventEntry := range em.events {
		event := eventEntry.event

		ok, _ := regexp.MatchString(event.Regex, eventPayload)
		if !ok {
			continue
		}

		r := event.Func(eventPayload)
		if r != nil {
			eventEntry.mu.Lock()
			eventEntry.queue = append(eventEntry.queue, r)
			eventEntry.mu.Unlock()

		}
	}
}

// Returns every event names
func (em *EconEventManager) Names() []string {
	ret := []string{}

	for name := range em.events {
		ret = append(ret, name)
	}

	return ret
}
