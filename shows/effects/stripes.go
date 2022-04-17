package effects

import (
	"image/color"

	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows/parameters"
)

// StripesEffect is a effect that draws moving stripes in one color
type StripesEffect struct {
	color  *parameters.Parameter
	speed  *parameters.Parameter
	length *parameters.Parameter
	gap    *parameters.Parameter

	currentPosition float64
}

// NewStripesEffect returns a new stripes effect
func NewStripesEffect() *StripesEffect {
	blink := StripesEffect{}

	blink.color = parameters.NewParameter("color", parameters.Color, "Color")
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
	color := e.color.Get().(color.NRGBA)
	speed := e.speed.Get().(int)
	length := e.length.Get().(int)
	gap := e.gap.Get().(int)

	numLeds := hw.Led.GetNumLedsMultiPart(parts)
	ledsPerSecond := mapPercent(0.0, 75.0, speed)
	pos := getNextPosition(&e.currentPosition, ledsPerSecond, numLeds, nanoseconds)

	// draw beginning from current position (we use the wrap around here)
	for i := 0; i < numLeds; i++ {
		if i%(length+gap) < length {
			hw.Led.SetColorMultiPart(parts, pos+i, color.R, color.G, color.B, true)
		} else {
			hw.Led.SetColorMultiPart(parts, pos+i, 0, 0, 0, true)
		}
	}
}

// Parameters returns the list of paremeters
func (e *StripesEffect) Parameters() []*parameters.Parameter {
	// todo: maybe only once?
	data := make([]*parameters.Parameter, 4)
	data[0] = e.color
	data[1] = e.speed
	data[2] = e.length
	data[3] = e.gap
	return data
}
