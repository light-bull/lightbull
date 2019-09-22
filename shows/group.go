package shows

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows/effects"
)

// Group maps an effect type to a group of LED parts
type Group struct {
	ID uuid.UUID `json:"id"`

	Effect effects.Effect `json:"-"`

	parts []string

	// FIXME: mux!
}

type groupJSON struct {
	ID     uuid.UUID           `json:"id"`
	Parts  []string            `json:"parts"`
	Effect *effects.EffectJSON `json:"effect"`
}

// newGroup creates a new group. It is meant to be called from Visual.
func newGroup(parts []string, effect string) (*Group, error) {
	group := Group{ID: uuid.New()} // FIXME: uuid is randomly generated, so there could be a collission

	err := group.SetParts(parts)
	if err != nil {
		return nil, err
	}

	err = group.SetEffect(effect)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

// MarshalJSON is there to implement the `json.Marshaller` interface.
func (group *Group) MarshalJSON() ([]byte, error) {
	data := groupJSON{ID: group.ID, Parts: group.parts}

	if group.Effect != nil {
		data.Effect = effects.EffectToJSON(group.Effect)
	}

	return json.Marshal(data)
}

// UnmarshalJSON is there to implement the `json.Unmarshaller` interface.
func (group *Group) UnmarshalJSON(data []byte) error {
	input := groupJSON{}

	err := json.Unmarshal(data, &input)
	if err != nil {
		return err
	}

	group.ID = input.ID
	group.parts = input.Parts

	effect := effects.EffectFromJSON(input.Effect)
	if effect != nil {
		group.Effect = *effect
	}

	// TODO: input validation

	return nil
}

// Parts returns the LED parts that are configured for this effect.
func (group *Group) Parts() []string {
	return group.parts
}

// SetParts changes the LED parts that are configured for this effect.
func (group *Group) SetParts(parts []string) error {
	group.parts = parts
	// TODO
	// TODO: check that part is only configured for one effect in a visual
	return nil
}

// SetEffect changes the effect type for this group
func (group *Group) SetEffect(effecttype string) error {
	effect := effects.NewEffect(effecttype)
	if effect == nil {
		return errors.New("Unknown effect")
	}

	group.Effect = effect
	return nil
}

// Update decides about the changes that are caused by the group/effect for a certain timestep.
func (group *Group) Update(hw *hardware.Hardware, nanoseconds int64) {
	if group.Effect != nil {
		group.Effect.Update(hw, group.parts, nanoseconds)
	}
}
