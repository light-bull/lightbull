package effects

import (
	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows/parameters"
)

type BlinkEffect struct {
	color *parameters.Parameter
	ratio *parameters.Parameter
	speed *parameters.Parameter
}

func NewBlinkEffect() *BlinkEffect {
	blink := BlinkEffect{}

	blink.color = parameters.NewParameter("color", parameters.Color, "Color")
	blink.ratio = parameters.NewParameter("ratio", parameters.UInt8, "Ratio")
	blink.speed = parameters.NewParameter("speed", parameters.UInt8, "Speed")

	return &blink
}

func (e *BlinkEffect) Type() string {
	return Blink
}

func (e *BlinkEffect) Name() string {
	return "Blink"
}

func (e *BlinkEffect) Update(hw *hardware.Hardware, parts []string, nanoseconds int64) {
	// TODO implement blink effect
	// must somehow keep state about the time progression between light and dark state
}

func (e *BlinkEffect) Parameters() []*parameters.Parameter {
	params := make([]*parameters.Parameter, 3)
	params[0] = e.color
	params[1] = e.ratio
	params[2] = e.speed
	return params
}