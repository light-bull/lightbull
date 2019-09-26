package hardware

import (
	"errors"

	"github.com/spf13/viper"
)

// Hardware controlls all connected hardware like LEDs, the ethernet interface or the controller board itself.
type Hardware struct {
	Led    *LED
	System *System
}

// New initializes the hardware
func New() (*Hardware, error) {
	hw := Hardware{}

	// create leds struct
	hw.Led = NewLED()

	// read led parts from config file
	ledsConfig := viper.Sub("leds")
	if ledsConfig == nil {
		return nil, errors.New("Missing LED part definition") // should not happen since it is added by default. but who knows ;)
	}

	type partFormat struct {
		Name string  `mapstructure:"name"`
		Leds [][]int `mapstructure:"leds"`
	}

	type format struct {
		Parts []partFormat `mapstructure:"parts"`
	}
	var partConfig format

	err := ledsConfig.Unmarshal(&partConfig)
	if err != nil {
		return nil, errors.New("Malformed LED part definition")
	}

	// configure led parts
	for _, part := range partConfig.Parts {
		for _, leds := range part.Leds {
			if len(leds) != 2 {
				return nil, errors.New("Malformed LED part definition")
			}

			hw.Led.AddPart(part.Name, leds[0], leds[1])
		}
	}

	// init leds
	if err := hw.Led.Init(); err != nil {
		return nil, err
	}

	// system
	hw.System = NewSystem()

	// finished
	return &hw, nil
}

// Update writes changes to the hardware
func (hw *Hardware) Update() {
	hw.Led.Update()
}
