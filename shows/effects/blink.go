package effects

import (
	"image/color"

	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows/parameters"
)

// BlinkEffect is a effect that lets the LEDs blink in one color
type BlinkEffect struct {
	color *parameters.Parameter
	speed *parameters.Parameter
	ratio *parameters.Parameter

	nsSinceLastStart int64
}

// NewBlinkEffect returns a new blink effect
func NewBlinkEffect() *BlinkEffect {
	blink := BlinkEffect{}

	blink.color = parameters.NewParameter("color", parameters.Color, "Color")
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
	color := e.color.Get().(color.NRGBA)
	speed := e.speed.Get().(int)
	ratio := e.ratio.Get().(int)

	// length of one on-off cycle
	interval := mapPercent(int64(100000000), 5000000000, speed)
	intervalOn := mapPercent(0, interval, ratio)

	// get time since last start of on-off cycle
	e.nsSinceLastStart = (e.nsSinceLastStart + nanoseconds) % interval

	// turn on or off
	var r, g, b byte = 0, 0, 0
	if e.nsSinceLastStart < intervalOn {
		r = color.R
		g = color.G
		b = color.B
	}

	for _, part := range parts {
		hw.Led.SetColorAllPart(part, r, g, b)
	}
}

// Parameters returns the list of paremeters
func (e *BlinkEffect) Parameters() []*parameters.Parameter {
	// todo: maybe only once?
	data := make([]*parameters.Parameter, 3)
	data[0] = e.color
	data[1] = e.speed
	data[2] = e.ratio
	return data
}
