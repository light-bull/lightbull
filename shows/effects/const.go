package effects

const (
	// Calibration is the calibration effect
	Calibration = "calibration"

	// SingleColor is the single color effect
	SingleColor = "singlecolor"

	// Blink is the blink effect
	Blink = "blink"

	// Stripes is the stripes effect
	Stripes = "stripes"

	// Rainbow is a rainbow effect
	Rainbow = "rainbow"
)

var effectNames map[string]string

// GetEffects returns the type and name of all effects
func GetEffects() map[string]string {
	// create map on first call of function, reuse later on
	if effectNames == nil {
		effectNames = make(map[string]string)

		effectNames[Calibration] = NewEffect(Calibration).Name()
		effectNames[SingleColor] = NewEffect(SingleColor).Name()
		effectNames[Blink] = NewEffect(Blink).Name()
		effectNames[Stripes] = NewEffect(Stripes).Name()
		effectNames[Rainbow] = NewEffect(Rainbow).Name()
	}

	return effectNames
}
