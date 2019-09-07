package effects

import "github.com/light-bull/lightbull/shows/parameters"

// SingleColor is a effect that lets the LEDs show one color
type SingleColor struct {
	color parameters.Color
}

// Type returns "singlecolor"
func (e *SingleColor) Type() string {
	return "singlecolor"
}

// Name returns "Single Color"
func (e *SingleColor) Name() string {
	return "Single Color"
}

// Update decides about the changes that are caused by the effect for a certain timestep.
func (e *SingleColor) Update(nanoseconds int64) {
	return
}

// Parameters returns the list of paremeters
func (e *SingleColor) Parameters() [](*parameters.Parameter) {
	return nil
}
