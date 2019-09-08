package effects

import (
	"image/color"

	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows/parameters"
)

// SingleColorEffect is a effect that lets the LEDs show one color
type SingleColorEffect struct {
	color *parameters.Parameter
}

// NewSingleColorEffect returns a new single color effect
func NewSingleColorEffect() *SingleColorEffect {
	singlecolor := SingleColorEffect{}

	singlecolor.color = parameters.NewParameter("color", parameters.Color, "Color")

	return &singlecolor
}

// Type returns "singlecolor"
func (e *SingleColorEffect) Type() string {
	return SingleColor
}

// Name returns "Single Color"
func (e *SingleColorEffect) Name() string {
	return "Single Color"
}

// Update decides about the changes that are caused by the effect for a certain timestep.
func (e *SingleColorEffect) Update(hw *hardware.Hardware, parts []string, nanoseconds int64) {
	color := e.color.Get().(color.NRGBA)

	for _, part := range parts {
		hw.Led.SetColorPart(part, color.R, color.G, color.B)
	}
}

// Parameters returns the list of paremeters
func (e *SingleColorEffect) Parameters() []*parameters.Parameter {
	// todo: maybe only once?
	data := make([]*parameters.Parameter, 1)
	data[0] = e.color
	return data
}
