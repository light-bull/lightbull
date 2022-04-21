package parameters

import (
	"encoding/json"
)

// this is a simple single-value parameter. see color.go for a more complex example

// BooleanType the datatype for boolean
type BooleanType struct {
	value bool
}

// NewBooleanType returns a new data of type BooleanType
func NewBooleanType() *BooleanType {
	boolean := BooleanType{}

	boolean.value = false

	return &boolean
}

// Type returns "boolean"
func (c *BooleanType) Type() string {
	return Boolean
}

// Get the value
func (c *BooleanType) Get() interface{} {
	return c.value
}

// Set the color
func (c *BooleanType) Set(new interface{}) error {
	var tmp = new.(bool)
	c.value = tmp
	return nil
}

// MarshalJSON returns the data serialized as JSON
func (c *BooleanType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

// UnmarshalJSON loads the data from the JSON string
func (c *BooleanType) UnmarshalJSON(data []byte) error {
	var input bool

	err := json.Unmarshal(data, &input)
	if err != nil {
		return err
	}

	return c.Set(input)
}
