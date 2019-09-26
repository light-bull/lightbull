package persistence

import (
	"github.com/light-bull/lightbull/events"
)

// EventClient handles event that should trigger a save operation
type EventClient struct {
	event chan *events.Event

	persistence *Persistence
}

// newEventClient create a new persistence event client
func newEventClient(persistence *Persistence) *EventClient {
	client := EventClient{
		persistence: persistence,
	}

	client.event = make(chan *events.Event)
	go client.run()

	return &client
}

// EventChan is there to implement the `EventClient` interface
func (client *EventClient) EventChan() chan *events.Event {
	return client.event
}

// run is the event handler loop
func (client *EventClient) run() {
	for {
		select {
		case event := <-client.event:
			if event.Show() != nil {
				switch event.Topic {
				case events.ShowDeleted:
					client.persistence.DeleteShow(event.Show())
				case events.ParameterChanged:
					// ignore, only changes to default values are written
				default:
					client.persistence.SaveShow(event.Show())
				}
			}
		}
	}
}
