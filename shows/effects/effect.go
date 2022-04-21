package effects

import (
	"encoding/json"
	"errors"

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

	// Parameters returns the list of parameters
	Parameters() []*parameters.Parameter
}

// NewEffect returns a new effect of specified effect type (or nil)
func NewEffect(effecttype string) Effect {
	if effecttype == SingleColor {
		return NewSingleColorEffect()
	} else if effecttype == Blink {
		return NewBlinkEffect()
	} else if effecttype == Stripes {
		return NewStripesEffect()
	} else if effecttype == Rainbow {
		return NewRainbowEffect()
	}
	return nil
}

// EffectJSON is the JSON format for effects
type EffectJSON struct {
	Type       string                  `json:"type"`
	Parameters []*parameters.Parameter `json:"parameters"`

	// only for deserialization
	effect *Effect
}

// EffectToJSON converts an effect and parameters to EffectJSON.
func EffectToJSON(effect Effect) *EffectJSON {
	output := EffectJSON{}

	output.Type = effect.Type()
	output.Parameters = effect.Parameters()

	return &output
}

// EffectFromJSON creates an effect (an object that fulfills the effect interface) from EffectJSON.
func EffectFromJSON(data *EffectJSON) *Effect {
	// The deserialization already happened in EffectJSON.UnmarshallJSON, lets just use the result
	return data.effect
}

// UnmarshalJSON deserializes an effect and parameter store
func (effectjson *EffectJSON) UnmarshalJSON(data []byte) error {
	// get the type first
	type format struct {
		Type       string             `json:"type"`
		Parameters []*json.RawMessage `json:"parameters"`
	}

	dataMap := format{}

	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		return err
	}

	// we know the type, so just create the corresponding effect
	effect := NewEffect(dataMap.Type)
	if effect == nil {
		return errors.New("invalid effect type")
	}

	// use the key to lookup the parameter in the effect. then call unmarshal on the parameter with the concrete datatypes
	for _, parameterRaw := range dataMap.Parameters {
		// parse the data of one parameter. we are only interested in the key
		type parameterFormat struct {
			Key string `json:"key"`
		}
		parameterMap := parameterFormat{}
		err := json.Unmarshal(*parameterRaw, &parameterMap)
		if err != nil {
			return err
		}

		// search for parameter
		for _, parameter := range effect.Parameters() {
			if parameter.Key == parameterMap.Key {
				err = parameter.UnmarshalJSON(*parameterRaw)
				if err != nil {
					return err
				}

				break
			}
		}
	}

	// we are done
	effectjson.effect = &effect
	return nil
}
