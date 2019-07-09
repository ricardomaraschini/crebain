package tui

// Event processor is an entity capable of processing user events.
type EventProcessor interface {
	Event(string)
}
