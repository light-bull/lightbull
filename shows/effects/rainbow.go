package effects

import (
	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows/parameters"
)

// RainbowEffect is a effect that draws a moving rainbow
type RainbowEffect struct {
	speed    *parameters.Parameter
	reversed *parameters.Parameter

	currentPosition float64
}

// NewRainbowEffect returns a new rainbow effect
func NewRainbowEffect() *RainbowEffect {
	rainbow := RainbowEffect{}

	rainbow.speed = parameters.NewParameter("speed", parameters.Percent, "Speed")
	rainbow.reversed = parameters.NewParameter("reversed", parameters.Boolean, "Reversed")

	return &rainbow
}

// Type returns "rainbow"
func (e *RainbowEffect) Type() string {
	return Rainbow
}

// Name returns "Rainbow"
func (e *RainbowEffect) Name() string {
	return "Rainbow"
}

// Update decides about the changes that are caused by the effect for a certain timestep.
func (e *RainbowEffect) Update(hw *hardware.Hardware, parts []string, nanoseconds int64) {
	speed := e.speed.Get().(int)
	reversed := e.reversed.Get().(bool)

	numLeds := hw.Led.GetNumLedsMultiPart(parts)
	ledsPerSecond := mapPercent(0.0, 300.0, speed)
	pos := getNextPosition(&e.currentPosition, ledsPerSecond, numLeds, nanoseconds, reversed)

	directionFactor := getDirectionFactor(reversed)

	for i := 0; i < numLeds; i++ {
		directionalIndex := i * directionFactor
		hue := moduloInt(directionalIndex*360/(numLeds-1), 360)
		r, g, b := hsv2rgb(hue, 100, 100)
		hw.Led.SetColorMultiPart(parts, pos+directionalIndex, r, g, b, true)
	}
}

// Parameters returns the list of paremeters
func (e *RainbowEffect) Parameters() []*parameters.Parameter {
	data := make([]*parameters.Parameter, 2)
	data[0] = e.speed
	data[1] = e.reversed
	return data
}
