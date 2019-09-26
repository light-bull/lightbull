package events

// EventClient is an interface that describes a client that can subscribe to events
type EventClient interface {
	EventChan() chan *Event
}
