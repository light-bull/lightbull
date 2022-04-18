package effects

import (
	"image/color"

	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows/parameters"
)

// StripesEffect is a effect that draws moving stripes in one color
type StripesEffect struct {
	color_primary   *parameters.Parameter
	color_secondary *parameters.Parameter
	speed           *parameters.Parameter
	length          *parameters.Parameter
	gap             *parameters.Parameter

	currentPosition float64
}

// NewStripesEffect returns a new stripes effect
func NewStripesEffect() *StripesEffect {
	blink := StripesEffect{}

	blink.color_primary = parameters.NewParameter("color_primary", parameters.Color, "Primary color")
	blink.color_secondary = parameters.NewParameter("color_secondary", parameters.Color, "Secondary color")
	blink.speed = parameters.NewParameter("speed", parameters.Percent, "Speed")
	blink.length = parameters.NewParameter("length", parameters.IntegerGreaterZero, "Length")
	blink.gap = parameters.NewParameter("gap", parameters.IntegerGreaterZero, "Gap")

	return &blink
}

// Type returns "blink"
func (e *StripesEffect) Type() string {
	return Stripes
}

// Name returns "Blink"
func (e *StripesEffect) Name() string {
	return "Stripes"
}

// Update decides about the changes that are caused by the effect for a certain timestep.
func (e *StripesEffect) Update(hw *hardware.Hardware, parts []string, nanoseconds int64) {
	color_primary := e.color_primary.Get().(color.NRGBA)
	color_secondary := e.color_secondary.Get().(color.NRGBA)
	speed := e.speed.Get().(int)
	length := e.length.Get().(int)
	gap := e.gap.Get().(int)

	numLeds := hw.Led.GetNumLedsMultiPart(parts)
	ledsPerSecond := mapPercent(0.0, 75.0, speed)
	pos := getNextPosition(&e.currentPosition, ledsPerSecond, numLeds, nanoseconds)

	// draw beginning from current position (we use the wrap around here)
	for i := 0; i < numLeds; i++ {
		if i%(length+gap) < length {
			hw.Led.SetColorMultiPart(parts, pos+i, color_primary.R, color_primary.G, color_primary.B, true)
		} else {
			hw.Led.SetColorMultiPart(parts, pos+i, color_secondary.R, color_secondary.G, color_secondary.B, true)
		}
	}
}

// Parameters returns the list of paremeters
func (e *StripesEffect) Parameters() []*parameters.Parameter {
	// todo: maybe only once?
	data := make([]*parameters.Parameter, 5)
	data[0] = e.color_primary
	data[1] = e.color_secondary
	data[2] = e.speed
	data[3] = e.length
	data[4] = e.gap
	return data
}
