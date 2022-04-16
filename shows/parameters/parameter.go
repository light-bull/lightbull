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

// NewParameter returns a new parameter of the specified data type (or nil)
func NewParameter(key string, datatype string, name string) *Parameter {
	parameter := Parameter{}

	parameter.ID = uuid.New() // TODO: make sure that unique
	parameter.Key = key
	parameter.Name = name

	if datatype == Color {
		parameter.cur = NewColor()
		parameter.def = NewColor()
	} else if datatype == Percent {
		parameter.cur = NewPercent()
		parameter.def = NewPercent()
	} else {
		return nil
	}

	return &parameter
}

// MarshalJSON is there to implement the `json.Marshaller` interface.
func (parameter *Parameter) MarshalJSON() ([]byte, error) {
	type format struct {
		ID      uuid.UUID `json:"id"`
		Key     string    `json:"key"`
		Name    string    `json:"name"` // will be ignored for deserialization
		Type    string    `json:"type"` // will be ignored for deserialization
		Default DataType  `json:"default"`
		Current DataType  `json:"current"`
	}

	data := format{
		ID:      parameter.ID,
		Key:     parameter.Key,
		Name:    parameter.Name,
		Type:    parameter.cur.Type(),
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
		err = parameter.SetFromJSON(*dataMap.Current)
		if err != nil {
			return err
		}
	}

	if dataMap.Default != nil {
		err = parameter.SetDefaultFromJSON(*dataMap.Default)
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

// SetFromJSON sets a new value from the JSON data
func (parameter *Parameter) SetFromJSON(data []byte) error {
	err := parameter.cur.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	return nil
}

// SetDefaultFromJSON sets a new default value from the JSON data
func (parameter *Parameter) SetDefaultFromJSON(data []byte) error {
	err := parameter.def.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	return nil
}

// SetDefault sets the current value as default
// TODO: remove?
func (parameter *Parameter) SetDefault() {
	parameter.def.Set(parameter.cur.Get())
}

// RestoreDefault sets the current value back to the default value
// TODO: remove?
func (parameter *Parameter) RestoreDefault() {
	parameter.cur.Set(parameter.def.Get())
}
