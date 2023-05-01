package parameters

import (
	"encoding/json"
	"errors"
)

// this is a simple single-value parameter. see color.go for a more complex example

// IntegerGreaterOrEqualZeroType the datatype for integers greater than 0
type IntegerGreaterOrEqualZeroType struct {
	value int
}

// NewIntegerGreaterZero returns a new data of type integergreaterorequalzero
func NewIntegerGreaterZero() *IntegerGreaterOrEqualZeroType {
	integer := IntegerGreaterOrEqualZeroType{}

	integer.value = 1

	return &integer
}

// Type returns "integergreaterorequalzero"
func (c *IntegerGreaterOrEqualZeroType) Type() string {
	return IntegerGreaterOrEqualZero
}

// Get the value
func (c *IntegerGreaterOrEqualZeroType) Get() interface{} {
	return c.value
}

// Set the color
func (c *IntegerGreaterOrEqualZeroType) Set(new interface{}) error {
	var tmp int = new.(int)
	if tmp < 0 {
		return errors.New("invalid value for parameter of type integergreaterorequalzero")
	}

	c.value = tmp
	return nil
}

// MarshalJSON returns the data serialized as JSON
func (c *IntegerGreaterOrEqualZeroType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

// UnmarshalJSON loads the data from the JSON string
func (c *IntegerGreaterOrEqualZeroType) UnmarshalJSON(data []byte) error {
	var input int

	err := json.Unmarshal(data, &input)
	if err != nil {
		return err
	}

	return c.Set(input)
}
