package effects

import (
	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows/parameters"
	"image/color"
)

// CalibrationEffect is an effect that sets a single LED to a color for calibration purposes
type CalibrationEffect struct {
	color *parameters.Parameter
	ledId *parameters.Parameter
}

// NewCalibrationEffect returns a new calibration effect
func NewCalibrationEffect() *CalibrationEffect {
	calibration := CalibrationEffect{}

	calibration.color = parameters.NewParameter("color", parameters.Color, "Color")
	calibration.ledId = parameters.NewParameter("ledId", parameters.IntegerGreaterOrEqualZero, "ID of the Led to set")

	return &calibration
}

// Type returns "calibration"
func (c *CalibrationEffect) Type() string {
	return Calibration
}

// Name returns "Calibration"
func (c *CalibrationEffect) Name() string {
	return "Calibration"
}

// Update decides about the changes that are caused by the effect for a certain timestep.
func (c *CalibrationEffect) Update(hw *hardware.Hardware, parts []string, nanoseconds int64) {
	primaryColor := c.color.Get().(color.NRGBA)
	ledId := c.ledId.Get().(int)

	for _, part := range parts {
		hw.Led.SetColorAllPart(part, 0, 0, 0)
		hw.Led.SetColor(part, ledId, primaryColor.R, primaryColor.G, primaryColor.B)
	}
}

// Parameters returns the list of parameters
func (c *CalibrationEffect) Parameters() []*parameters.Parameter {
	data := make([]*parameters.Parameter, 2)
	data[0] = c.color
	data[1] = c.ledId
	return data
}
