package effects

const (
	// SingleColor is the single color effect
	SingleColor = "singlecolor"
)

var effectNames map[string]string

// GetEffects returns the type and name of all effects
func GetEffects() map[string]string {
	// create list on first call of function, reuse lateron
	if effectNames == nil {
		effectNames = make(map[string]string)

		effectNames[SingleColor] = NewEffect(SingleColor).Name()
	}

	return effectNames
}
