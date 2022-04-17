package effects

const (
	// SingleColor is the single color effect
	SingleColor = "singlecolor"

	// Blink is the blink effect
	Blink = "blink"

	// Stripes is the stripes effect
	Stripes = "stripes"
)

var effectNames map[string]string

// GetEffects returns the type and name of all effects
func GetEffects() map[string]string {
	// create map on first call of function, reuse later on
	if effectNames == nil {
		effectNames = make(map[string]string)

		effectNames[SingleColor] = NewEffect(SingleColor).Name()
		effectNames[Blink] = NewEffect(Blink).Name()
		effectNames[Stripes] = NewEffect(Stripes).Name()
	}

	return effectNames
}
