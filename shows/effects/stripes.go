package effects

import (
	"image/color"

	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows/parameters"
)

// StripesEffect is a effect that draws moving stripes in one color
type StripesEffect struct {
	colorPrimary   *parameters.Parameter
	colorSecondary *parameters.Parameter
	speed          *parameters.Parameter
	length         *parameters.Parameter
	gap            *parameters.Parameter
	reversed       *parameters.Parameter

	currentPosition float64
}

// NewStripesEffect returns a new stripes effect
func NewStripesEffect() *StripesEffect {
	stripes := StripesEffect{}

	stripes.colorPrimary = parameters.NewParameter("colorPrimary", parameters.Color, "Primary color")
	stripes.colorSecondary = parameters.NewParameter("colorSecondary", parameters.Color, "Secondary color")
	stripes.speed = parameters.NewParameter("speed", parameters.Percent, "Speed")
	stripes.length = parameters.NewParameter("length", parameters.IntegerGreaterZero, "Length")
	stripes.gap = parameters.NewParameter("gap", parameters.IntegerGreaterZero, "Gap")
	stripes.reversed = parameters.NewParameter("reversed", parameters.Boolean, "Reversed")

	return &stripes
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
	colorPrimary := e.colorPrimary.Get().(color.NRGBA)
	colorSecondary := e.colorSecondary.Get().(color.NRGBA)
	speed := e.speed.Get().(int)
	length := e.length.Get().(int)
	gap := e.gap.Get().(int)
	reversed := e.reversed.Get().(bool)

	numLeds := hw.Led.GetNumLedsMultiPart(parts)
	ledsPerSecond := mapPercent(0.0, 75.0, speed)
	pos := getNextPosition(&e.currentPosition, ledsPerSecond, numLeds, nanoseconds, reversed)

	directionFactor := getDirectionFactor(reversed)

	// draw beginning from current position (we use the wrap around here)
	for i := 0; i < numLeds; i++ {
		if i%(length+gap) < length {
			hw.Led.SetColorMultiPart(parts, pos+directionFactor*i, colorPrimary.R, colorPrimary.G, colorPrimary.B, true)
		} else {
			hw.Led.SetColorMultiPart(parts, pos+directionFactor*i, colorSecondary.R, colorSecondary.G, colorSecondary.B, true)
		}
	}
}

// Parameters returns the list of paremeters
func (e *StripesEffect) Parameters() []*parameters.Parameter {
	data := make([]*parameters.Parameter, 6)
	data[0] = e.colorPrimary
	data[1] = e.colorSecondary
	data[2] = e.speed
	data[3] = e.length
	data[4] = e.gap
	data[5] = e.reversed
	return data
}
