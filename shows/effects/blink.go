package effects

import (
	"image/color"

	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows/parameters"
)

// BlinkEffect is a effect that lets the LEDs blink in one color
type BlinkEffect struct {
	color_primary   *parameters.Parameter
	color_secondary *parameters.Parameter
	speed           *parameters.Parameter
	ratio           *parameters.Parameter

	nsSinceLastStart int64
}

// NewBlinkEffect returns a new blink effect
func NewBlinkEffect() *BlinkEffect {
	blink := BlinkEffect{}

	blink.color_primary = parameters.NewParameter("color_primary", parameters.Color, "Primary color")
	blink.color_secondary = parameters.NewParameter("color_secondary", parameters.Color, "Secondary color")
	blink.speed = parameters.NewParameter("speed", parameters.Percent, "Speed")
	blink.ratio = parameters.NewParameter("ratio", parameters.Percent, "Ratio")

	return &blink
}

// Type returns "blink"
func (e *BlinkEffect) Type() string {
	return Blink
}

// Name returns "Blink"
func (e *BlinkEffect) Name() string {
	return "Blink"
}

// Update decides about the changes that are caused by the effect for a certain timestep.
func (e *BlinkEffect) Update(hw *hardware.Hardware, parts []string, nanoseconds int64) {
	color_primary := e.color_primary.Get().(color.NRGBA)
	color_secondary := e.color_secondary.Get().(color.NRGBA)
	speed := e.speed.Get().(int)
	ratio := e.ratio.Get().(int)

	// length of one on-off cycle
	interval := mapPercent(int64(5000000000), 100000000, speed)
	intervalOn := mapPercent(0, interval, ratio)

	// get time since last start of on-off cycle
	e.nsSinceLastStart = (e.nsSinceLastStart + nanoseconds) % interval

	// turn on or off
	var r, g, b byte = 0, 0, 0
	if e.nsSinceLastStart < intervalOn {
		r = color_primary.R
		g = color_primary.G
		b = color_primary.B
	} else {
		r = color_secondary.R
		g = color_secondary.G
		b = color_secondary.B
	}

	for _, part := range parts {
		hw.Led.SetColorAllPart(part, r, g, b)
	}
}

// Parameters returns the list of paremeters
func (e *BlinkEffect) Parameters() []*parameters.Parameter {
	// todo: maybe only once?
	data := make([]*parameters.Parameter, 4)
	data[0] = e.color_primary
	data[1] = e.color_secondary
	data[2] = e.speed
	data[3] = e.ratio
	return data
}
