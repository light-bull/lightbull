package events

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/light-bull/lightbull/shows"
)

const (
	// ShowAdded is the event topic when a new show was added
	ShowAdded = "show_added"

	// ShowChanged is the event topic when properties of a show itself were changed
	ShowChanged = "show_changed"

	// ShowDeleted is the event topic when a show was deleted
	ShowDeleted = "show_deleted"

	// VisualAdded is the event topic when a new visual was added
	VisualAdded = "visual_added"

	// VisualChanged is the event topic when properties of a visual itself were changed
	VisualChanged = "visual_changed"

	// VisualDeleted is the event topic when a visual was deleted
	VisualDeleted = "visual_deleted"

	// GroupAdded is the event topic when a new group was added
	GroupAdded = "group_added"

	// GroupChanged is the event topic when the group or the associated effect were changed (not for properties)
	GroupChanged = "group_changed"

	// GroupDeleted is the event topic when a group was deleted
	GroupDeleted = "group_deleted"

	// ParameterChanged is the event topic when the current value of a parameter was changed
	ParameterChanged = "parameter_changed"

	// ParameterDefaultChanged is the event topic when at least the default value of a parameter was changed. This can also include a change of the current value.
	ParameterDefaultChanged = "parameter_default_changed"

	// ParameterLinksChanged is the event topic when the links of a parameter changed
	ParameterLinksChanged = "parameter_links_changed"

	// CurrentChanged is the event topic when the current show or visual were changed
	CurrentChanged = "current_changed"
)

// EventMetaInfo stores meta information about the event
type EventMetaInfo struct {
	// ConnectionID is the client id that triggered this event
	ConnectionID uuid.UUID
}

// MarshalJSON is there to implement the `json.Marshaller` interface.
func (meta *EventMetaInfo) MarshalJSON() ([]byte, error) {
	type format struct {
		ConnectionID string `json:"connectionId,omitempty"`
	}
	data := format{}

	if meta.ConnectionID != uuid.Nil {
		data.ConnectionID = meta.ConnectionID.String()
	}

	return json.Marshal(data)
}

// Event is an event that is sent through the event hub
type Event struct {
	// Topic of the event (like show_added)
	Topic string `json:"topic"`

	// The payload, usually this is the changed objects
	Payload interface{} `json:"payload"`

	// Meta info about the event (which client triggered it?)
	Meta *EventMetaInfo `json:"meta"`

	// show is the show where the event belongs to
	show *shows.Show
}

// NewEvent creates a new event
func NewEvent(topic string, payload interface{}, show *shows.Show, connectionID uuid.UUID) *Event {
	event := Event{
		Topic:   topic,
		Payload: payload,
		Meta: &EventMetaInfo{
			ConnectionID: connectionID,
		},
		show: show,
	}

	return &event
}

// Show returns the associated show or `nil`
func (event *Event) Show() *shows.Show {
	return event.show
}
