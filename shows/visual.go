package shows

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"

	"github.com/google/uuid"
	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows/parameters"
)

// Visual is a collection of effects that are applied to LED parts and bundled parameters.
type Visual struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	groups []*Group

	mux sync.Mutex
}

// showJSON is the format for a serialized JSON configuration
type visualJSON struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Groups []*Group  `json:"groups"`
}

// newVisual creates a new visual. It is meant to be called from Show.
func newVisual(name string) *Visual {
	visual := Visual{ID: uuid.New(), Name: name} // FIXME: uuid is randomly generated, so there could be a collission

	return &visual
}

// MarshalJSON is there to implement the `json.Marshaller` interface.
func (visual *Visual) MarshalJSON() ([]byte, error) {
	data := visualJSON{ID: visual.ID, Name: visual.Name, Groups: visual.groups}
	return json.Marshal(data)
}

// UnmarshalJSON is there to implement the `json.Unmarshaller` interface.
func (visual *Visual) UnmarshalJSON(data []byte) error {
	input := visualJSON{}

	err := json.Unmarshal(data, &input)
	if err != nil {
		return err
	}

	visual.ID = input.ID
	visual.Name = input.Name
	visual.groups = input.Groups

	// TODO: input validation

	return nil
}

// Groups returns the list of groups.
func (visual *Visual) Groups() []*Group {
	return visual.groups
}

// NewGroup adds a new group with an effect to the visual.
func (visual *Visual) NewGroup(parts []string, effect string) (*Group, error) {
	group, err := newGroup(parts, effect)
	if err != nil {
		return nil, err
	}

	visual.mux.Lock()
	visual.groups = append(visual.groups, group)
	visual.mux.Unlock()

	return group, nil
}

// DeleteGroup adds a new group with an effect to the visual.
func (visual *Visual) DeleteGroup(group *Group) {
	visual.mux.Lock()
	defer visual.mux.Unlock()

	for pos, cur := range visual.groups {
		if group.ID == cur.ID {
			visual.groups = append(visual.groups[:pos], visual.groups[pos+1:]...)
			break
		}
	}
}

// Update decides about the changes that are caused by the visual for a certain timestep.
func (visual *Visual) Update(hw *hardware.Hardware, nanoseconds int64) {
	for _, group := range visual.groups {
		group.Update(hw, nanoseconds)
	}
}

// FindParameter returns the parameter with the given ID and the belonging group or nil for malformed and non-existing IDs
func (visual *Visual) FindParameter(id uuid.UUID) (*Group, *parameters.Parameter) {
	visual.mux.Lock()
	defer visual.mux.Unlock()

	// iterate over shows, visuals and groups
	for _, group := range visual.Groups() {
		for _, parameter := range group.Effect.Parameters() {
			if parameter.ID == id {
				return group, parameter
			}
		}
	}

	return nil, nil
}

// LinkParameter creates a link between two parameters
func (visual *Visual) LinkParameter(parameter1 *parameters.Parameter, parameter2 *parameters.Parameter) error {
	if reflect.TypeOf(parameter1.Get()) != reflect.TypeOf(parameter2.Get()) {
		return errors.New("cannot create link between parameters of different type")
	}

	if parameter1.ID == parameter2.ID {
		return errors.New("cannot create link between identical parameters")
	}

	// TODO: check that parameters belong to same visual

	visual.mux.Lock()
	defer visual.mux.Unlock()

	parameter1.AddLink(parameter2)
	parameter2.AddLink(parameter1)

	return nil
}

// UnlinkParameter deletes a link between two parameters
func (visual *Visual) UnlinkParameter(parameter1 *parameters.Parameter, parameter2 *parameters.Parameter) error {
	visual.mux.Lock()
	defer visual.mux.Unlock()

	parameter1.DeleteLink(parameter2)
	parameter2.DeleteLink(parameter1)

	return nil
}
