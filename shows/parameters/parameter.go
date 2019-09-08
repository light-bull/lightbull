package parameters

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Parameter is an effect parameter
type Parameter struct {
	// ID is the globally unique UUID for this parameter
	ID uuid.UUID

	// Key is the id that is unique for a single effect
	Key string

	// Name is the nice name for the UI
	Name string

	// current and default value of paramter, they have to be the same DataType
	cur DataType
	def DataType
}

type parameterJSON struct {
	ID      uuid.UUID `json:"id"`
	Key     string    `json:"key"`
	Default DataType  `json:"default"`
	Current DataType  `json:"current"`
}

// NewParameter returns a new parameter of the specified data type (or nil)
func NewParameter(key string, datatype string, name string) *Parameter {
	parameter := Parameter{}

	parameter.ID = uuid.New() // TODO: make sure that unique
	parameter.Key = key
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
		Key:     parameter.Key,
		Current: parameter.cur,
		Default: parameter.def,
	}
	return json.Marshal(data)
}

// UnmarshalJSON is there to implement the `json.Unmarshaller` interface.
func (parameter *Parameter) UnmarshalJSON(data []byte) error {
	type format struct {
		ID      uuid.UUID        `json:"id"`
		Key     string           `json:"key"`
		Current *json.RawMessage `json:"current"`
		Default *json.RawMessage `json:"default"`
	}

	dataMap := format{}

	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		return err
	}

	parameter.ID = dataMap.ID
	parameter.Key = dataMap.Key

	if dataMap.Current != nil {
		err = parameter.cur.UnmarshalJSON(*dataMap.Current)
		if err != nil {
			return err
		}
	}

	if dataMap.Default != nil {
		err = parameter.def.UnmarshalJSON(*dataMap.Default)
		if err != nil {
			return err
		}
	}

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
