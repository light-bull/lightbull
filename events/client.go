package events

import "fmt"

// EventClient is an interface that describes a client that can subscribe to events
type EventClient interface {
	SendChan() chan *Event
}

// TODO: remove me
type EventDebugClient struct {
	send chan *Event
}

func NewEventDebugClient() *EventDebugClient {
	client := EventDebugClient{}

	client.send = make(chan *Event)
	go client.run()

	return &client
}

func (client *EventDebugClient) SendChan() chan *Event {
	return client.send
}

func (client *EventDebugClient) run() {
	for {
		select {
		case event := <-client.send:
			fmt.Println(event.Topic)
			fmt.Println(event.Payload)
		}
	}
}
