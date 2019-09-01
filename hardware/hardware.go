package hardware

// Hardware controlls all connected hardware like LEDs, the ethernet interface or the controller board itself.
type Hardware struct {
	Led    *LED
	System *System
}

// New initializes the hardware
func New() (*Hardware, error) {
	hw := Hardware{}

	// initialize leds
	hw.Led = NewLED()

	// TODO: move to config file
	hw.Led.AddPart("horn_left", 0, 68)
	hw.Led.AddPart("head_left", 69, 156)
	hw.Led.AddPart("head_left", 199, 249)
	hw.Led.AddPart("hole_left", 157, 198)
	hw.Led.AddPart("head_right", 250, 392)
	hw.Led.AddPart("horn_right", 400, 468)

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
