package parameters

import (
	"encoding/json"
	"image/color"
	"sync"
)

// ColorType is a datatype for RGB values
type ColorType struct {
	value color.NRGBA

	mux sync.Mutex
}

type colorDataJSON struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
}

// NewColor returns a new data of type color
func NewColor() *ColorType {
	color := ColorType{}

	color.value.G = 255 // FIXME: remove

	return &color
}

// Type returns "color"
func (c *ColorType) Type() string {
	return "color"
}

// Get the color
func (c *ColorType) Get() interface{} {
	return c.value
}

// Set the color
func (c *ColorType) Set(new interface{}) {
	c.value = new.(color.NRGBA)
}

// MarshalJSON returns the data serialized as JSON
func (c *ColorType) MarshalJSON() ([]byte, error) {
	c.mux.Lock()
	data := colorDataJSON{
		R: c.value.R,
		G: c.value.G,
		B: c.value.B,
	}
	c.mux.Unlock()

	return json.Marshal(data)
}

// UnmarshalJSON loads the data from the JSON string
func (c *ColorType) UnmarshalJSON(data []byte) error {
	input := colorDataJSON{}

	err := json.Unmarshal(data, &input)
	if err != nil {
		return err
	}

	c.mux.Lock()
	c.value.R = input.R
	c.value.G = input.G
	c.value.B = input.B
	c.mux.Unlock()

	return nil
}
