package teeworldsecon

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
}
