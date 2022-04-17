package parameters

import (
	"encoding/json"
	"errors"
)

// this is a simple single-value parameter. see color.go for a more complex example

// IntegerGreaterZeroType the datatype for integers greater than 0
type IntegerGreaterZeroType struct {
	value int
}

// NewIntegerGreaterZero returns a new data of type integergreaterzero
func NewIntegerGreaterZero() *IntegerGreaterZeroType {
	integer := IntegerGreaterZeroType{}

	integer.value = 1

	return &integer
}

// Type returns "integergreaterzero"
func (c *IntegerGreaterZeroType) Type() string {
	return IntegerGreaterZero
}

// Get the value
func (c *IntegerGreaterZeroType) Get() interface{} {
	return c.value
}

// Set the color
func (c *IntegerGreaterZeroType) Set(new interface{}) error {
	var tmp int = new.(int)
	if tmp < 0 {
		return errors.New("invalid value for parameter of type integergreaterzero")
	}

	c.value = tmp
	return nil
}

// MarshalJSON returns the data serialized as JSON
func (c *IntegerGreaterZeroType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

// UnmarshalJSON loads the data from the JSON string
func (c *IntegerGreaterZeroType) UnmarshalJSON(data []byte) error {
	var input int

	err := json.Unmarshal(data, &input)
	if err != nil {
		return err
	}

	return c.Set(input)
}
