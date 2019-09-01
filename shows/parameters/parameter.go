package parameters

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Parameter is an effect parameter
type Parameter struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	// current and default value of paramter, they have to be the same DataType
	cur DataType
	def DataType
}

type parameterJSON struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Current []byte    `json:"current"`
	Default []byte    `json:"default"`
}

// MarshalJSON is there to implement the `json.Marshaller` interface.
func (parameter *Parameter) MarshalJSON() ([]byte, error) {
	data := parameterJSON{ID: parameter.ID, Name: parameter.Name, Current: parameter.cur.ToJSON(), Default: parameter.def.ToJSON()}
	return json.Marshal(data)
}

// UnmarshalJSON is there to implement the `json.Unmarshaller` interface.
func (parameter *Parameter) UnmarshalJSON(data []byte) error {
	input := parameterJSON{}

	err := json.Unmarshal(data, &input)
	if err != nil {
		return err
	}

	parameter.ID = input.ID
	parameter.Name = input.Name
	// FIXME: set cur and def

	// TODO: input validation

	return nil
}

// Get returns the currently set value
func (parameter *Parameter) Get() interface{} {
	return parameter.cur.Get()
}

// ToJSON returns the currently set value as JSON
func (parameter *Parameter) ToJSON() []byte {
	return parameter.cur.ToJSON()
}

// UpdateFromJSON sets a new value from the JSON data
func (parameter *Parameter) UpdateFromJSON(data []byte) error {
	return parameter.cur.UpdateFromJSON(data)
}

// SetDefault sets the current value as default
func (parameter *Parameter) SetDefault() {
	parameter.def.Set(parameter.cur.Get())
}

// RestoreDefault sets the current value back to the default value
func (parameter *Parameter) RestoreDefault() {
	parameter.cur.Set(parameter.def.Get())
}
