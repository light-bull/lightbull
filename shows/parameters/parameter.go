package parameters

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Parameter is an effect parameter
type Parameter struct {
	ID   uuid.UUID
	Name string

	// current and default value of paramter, they have to be the same DataType
	cur DataType
	def DataType
}

type parameterJSON struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Current DataType  `json:"current"`
	Default DataType  `json:"default"`
}

// NewParameter returns a new parameter of the specified data type (or nil)
func NewParameter(name string, datatype string) *Parameter {
	parameter := Parameter{}

	parameter.ID = uuid.New() // TODO: make sure that unique
	parameter.Name = name

	if datatype == Color {
		parameter.cur = NewColor()
		parameter.def = NewColor()
	} else {
		return nil
	}

	return &parameter
}

// MarshalJSON is there to implement the `json.Marshaller` interface.
func (parameter *Parameter) MarshalJSON() ([]byte, error) {
	data := parameterJSON{
		ID:      parameter.ID,
		Name:    parameter.Name,
		Current: parameter.cur,
		Default: parameter.def,
	}
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

/*
// ToJSON returns the currently set value as JSON
func (parameter *Parameter) ToJSON() []byte {
	// TODO: error handling?
	data, _ := parameter.cur.MarshalJSON()
	return data
}

// UpdateFromJSON sets a new value from the JSON data
func (parameter *Parameter) UpdateFromJSON(data []byte) error {
	return parameter.cur.UpdateFromJSON(data)
}
*/

// SetDefault sets the current value as default
func (parameter *Parameter) SetDefault() {
	parameter.def.Set(parameter.cur.Get())
}

// RestoreDefault sets the current value back to the default value
func (parameter *Parameter) RestoreDefault() {
	parameter.cur.Set(parameter.def.Get())
}
