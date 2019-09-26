package events

import (
	"github.com/google/uuid"
	"github.com/light-bull/lightbull/shows"
)

// EventHub receives change events and distributes those to the subscribed clients
type EventHub struct {
	clients map[EventClient]bool

	// handle clients
	register   chan EventClient
	unregister chan EventClient

	// publish new events
	publish chan *Event
}

// NewEventHub returns a new event hub
func NewEventHub() *EventHub {
	eventhub := EventHub{
		clients:    make(map[EventClient]bool),
		register:   make(chan EventClient),
		unregister: make(chan EventClient),
		publish:    make(chan *Event),
	}

	go eventhub.run()

	return &eventhub
}

// RegisterClient registeres a new client at the event hub
func (eventhub *EventHub) RegisterClient(client EventClient) {
	eventhub.register <- client
}

// UnregisterClient unregisteres a client from the event hub
func (eventhub *EventHub) UnregisterClient(client EventClient) {
	eventhub.unregister <- client
}

// Publish a new event
func (eventhub *EventHub) Publish(event *Event) {
	eventhub.publish <- event
}

// PublishNew creates a new event and published it
func (eventhub *EventHub) PublishNew(topic string, payload interface{}, associatedShow *shows.Show, connectionID uuid.UUID) {
	// if the payload is a show, serialize only the shortened data
	show, ok := payload.(*shows.Show)
	if ok {
		payload = show.GetData()
	}

	// create event and publish it
	event := NewEvent(topic, payload, associatedShow, connectionID)
	eventhub.Publish(event)
}

// run startes the event hub so that events are distributed and clients are handled
func (eventhub *EventHub) run() {
	for {
		select {
		case client := <-eventhub.register:
			eventhub.clients[client] = true
		case client := <-eventhub.unregister:
			if _, ok := eventhub.clients[client]; ok {
				delete(eventhub.clients, client)
				close(client.EventChan())
			}
		case event := <-eventhub.publish:
			for client := range eventhub.clients {
				client.EventChan() <- event
				// TODO?
				//select {
				//case client.EventChan() <- event:
				//default:
				//	close(client.EventChan())
				//	delete(eventhub.clients, client)
				//}
			}
		}
	}
}
