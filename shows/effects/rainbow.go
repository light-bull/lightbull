package effects

import (
	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows/parameters"
)

// RainbowEffect is a effect that draws a moving rainbow
type RainbowEffect struct {
	speed *parameters.Parameter

	currentPosition float64
}

// NewRainbowEffect returns a new rainbow effect
func NewRainbowEffect() *RainbowEffect {
	rainbow := RainbowEffect{}

	rainbow.speed = parameters.NewParameter("speed", parameters.Percent, "Speed")

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

	numLeds := hw.Led.GetNumLedsMultiPart(parts)
	ledsPerSecond := mapPercent(0.0, 300.0, speed)
	pos := getNextPosition(&e.currentPosition, ledsPerSecond, numLeds, nanoseconds, false)

	for i := 0; i < numLeds; i++ {
		hue := (i * 360 / (numLeds - 1)) % 360
		r, g, b := hsv2rgb(hue, 100, 100)
		hw.Led.SetColorMultiPart(parts, i+pos, r, g, b, true)
	}
}

// Parameters returns the list of paremeters
func (e *RainbowEffect) Parameters() []*parameters.Parameter {
	// todo: maybe only once?
	data := make([]*parameters.Parameter, 1)
	data[0] = e.speed
	return data
}
