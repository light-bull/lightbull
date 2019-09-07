package parameters

import (
	"encoding/json"
	"image/color"
	"sync"
)

// Color is a datatype for RGB values
type Color struct {
	value color.NRGBA

	mux sync.Mutex
}

type colorDataJSON struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
}

// Type returns "color"
func (c *Color) Type() string {
	return "color"
}

// Get the color
func (c *Color) Get() interface{} {
	return c.value
}

// Set the color
func (c *Color) Set(new interface{}) {
	c.value = new.(color.NRGBA)
}

// MarshalJSON returns the data serialized as JSON
func (c *Color) MarshalJSON() ([]byte, error) {
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
func (c *Color) UnmarshalJSON(data []byte) error {
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
