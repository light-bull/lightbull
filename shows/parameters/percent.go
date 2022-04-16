package parameters

import (
	"encoding/json"
	"errors"
)

// this is a simple single-value parameter. see color.go for a more complex example

// PercentType is a datatype for 0 - 100%
type PercentType struct {
	value int
}

// NewPercent returns a new data of type percent
func NewPercent() *PercentType {
	percent := PercentType{}

	percent.value = 100

	return &percent
}

// Type returns "percent"
func (c *PercentType) Type() string {
	return Percent
}

// Get the value
func (c *PercentType) Get() interface{} {
	return c.value
}

// Set the color
func (c *PercentType) Set(new interface{}) error {
	var tmp int = new.(int)
	if tmp < 0 || tmp > 100 {
		return errors.New("invalid value for parameter of type percent")
	}

	c.value = tmp
	return nil
}

// MarshalJSON returns the data serialized as JSON
func (c *PercentType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

// UnmarshalJSON loads the data from the JSON string
func (c *PercentType) UnmarshalJSON(data []byte) error {
	var input int

	err := json.Unmarshal(data, &input)
	if err != nil {
		return err
	}

	return c.Set(input)
}
