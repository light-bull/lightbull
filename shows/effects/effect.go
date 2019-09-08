package effects

import (
	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows/parameters"
)

// Effect is the interface for effect implentations (like blink, sripes, ...)
type Effect interface {
	// Type returns a identifier like "blink" or "burn"
	Type() string

	// Name returns a nice name like "Blink"
	Name() string

	// Update decides about the changes that are caused by the effect for a certain timestep.
	Update(hw *hardware.Hardware, parts []string, nanoseconds int64)

	// Parameters returns the list of paremeters
	Parameters() [](*parameters.Parameter)
}

// NewEffect returns a new effect of specified effect type (or nil)
func NewEffect(effecttype string) Effect {
	if effecttype == SingleColor {
		return NewSingleColorEffect()
	}
	return nil
}

// EffectJSON is the JSON format for effects
type EffectJSON struct {
	Type       string                    `json:"type"`
	Name       string                    `json:"name"`
	Parameters [](*parameters.Parameter) `json:"parameters"`
}

// EffectToJSON converts an effect and parameters to EffectJSON.
func EffectToJSON(effect Effect) *EffectJSON {
	output := EffectJSON{}

	output.Type = effect.Type()
	output.Name = effect.Name()
	output.Parameters = effect.Parameters()

	return &output
}

// EffectFromJSON creates an effect (an object that fulfills the effect interface) from EffectJSON.
func EffectFromJSON(data *EffectJSON) *Effect {
	// TODO
	return nil
}
