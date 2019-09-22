package events

import "fmt"

// EventClient is an interface that describes a client that can subscribe to events
type EventClient interface {
	EventChan() chan *Event
}

// TODO: remove me
type EventDebugClient struct {
	event chan *Event
}

func NewEventDebugClient() *EventDebugClient {
	client := EventDebugClient{}

	client.event = make(chan *Event)
	go client.run()

	return &client
}

func (client *EventDebugClient) EventChan() chan *Event {
	return client.event
}

func (client *EventDebugClient) run() {
	for {
		select {
		case event := <-client.event:
			fmt.Println(event.Topic)
			fmt.Println(event.Payload)
		}
	}
}
