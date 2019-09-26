package hardware

import (
	"errors"
	"image"
	"image/color"
	"log"

	"github.com/spf13/viper"

	"periph.io/x/extra/devices/screen"
	"periph.io/x/periph/conn/display"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/devices/apa102"
	"periph.io/x/periph/host"
)

// LED is used to interact with the LED stripes. First, add the single parts and then run Init.
type LED struct {
	spi         spi.PortCloser
	apa102      display.Drawer
	apa102Dummy bool
	image       *image.NRGBA

	parts      []string
	partLedMap map[string][]int
	maxLedID   int

	maxColorSum int
	drawDummy   bool
}

// NewLED creates a new LED struct. After that, `AddPart` needs to be called and then `Init`.
func NewLED() *LED {
	led := &LED{}
	led.partLedMap = make(map[string][]int)
	led.maxLedID = -1
	return led
}

// AddPart adds a part definition for the LED stripes. The LEDs from id `first` to id `last`
// are added to the part `name`. AddPart can be called multiple times per part.
func (led *LED) AddPart(name string, first int, last int) {
	// add to ordered list of part names if necessary
	if !led.HasPart(name) {
		led.parts = append(led.parts, name)
	}

	// add to mapping
	led.partLedMap[name] = append(led.partLedMap[name], getRange(first, last)...)

	// update max led id
	if last > led.maxLedID {
		led.maxLedID = last
	}
}

// Init connects to the LED stripes
func (led *LED) Init() error {
	// check number of LEDs
	numLeds := led.getTotalNumLeds()
	if numLeds == 0 {
		return errors.New("No LEDs defined")
	}

	// initialize periph library
	if _, err := host.Init(); err != nil {
		return err
	}

	// initialize spispi.spiConn
	spiConn, err := spireg.Open("")
	if err != nil {
		log.Print("Failed to find a SPI port, printing at the console:\n")
		led.apa102 = screen.New(numLeds)
		led.apa102Dummy = true
	} else {
		// SPI config
		led.spi = spiConn
		led.spi.LimitSpeed(physic.Frequency(viper.GetInt("leds.spiMHz")) * physic.MegaHertz)

		// initialize apa102
		opts := apa102.DefaultOpts
		opts.NumPixels = numLeds
		apa102Dev, err := apa102.New(spiConn, &opts)
		if err != nil {
			return err
		}
		led.apa102 = apa102Dev
		led.apa102Dummy = false
	}

	// initialize image memory
	led.image = image.NewNRGBA(led.apa102.Bounds())

	// set brightness cap
	led.maxColorSum = (3 * 255) * viper.GetInt("leds.brightnessCap") / 100

	// load other config
	led.drawDummy = viper.GetBool("leds.drawDummy")

	return nil
}

// GetParts returns the names of all parts
func (led *LED) GetParts() []string {
	return led.parts
}

// HasPart checks if `part` is a valid part name
func (led *LED) HasPart(part string) bool {
	_, exists := led.partLedMap[part]
	return exists
}

// GetNumLeds returns the number of leds in a part
func (led *LED) GetNumLeds(part string) int {
	ledIDs, exists := led.partLedMap[part]
	if !exists {
		panic("invalid part name")
	}
	return len(ledIDs)
}

// SetColor sets the color for one pixel. UpdateColors needs to be called to make the changes visible
func (led *LED) SetColor(part string, pos int, r byte, g byte, b byte) {
	ledID := led.mapLedPartPos(part, pos)

	// filter colors that would need to much power
	sum := int(r) + int(g) + int(b)
	if sum > led.maxColorSum {
		diff := sum - led.maxColorSum
		r -= byte(diff * int(r) / sum)
		g -= byte(diff * int(g) / sum)
		b -= byte(diff * int(b) / sum)
	}

	led.image.SetNRGBA(ledID, 0, color.NRGBA{R: r, G: g, B: b, A: 255})
}

// SetColorPart sets the color for a whole part. UpdateColors needs to be called to make the changes visible
func (led *LED) SetColorPart(part string, r byte, g byte, b byte) {
	for i := 0; i < led.GetNumLeds(part); i++ {
		led.SetColor(part, i, r, g, b)
	}
}

// SetColorAll sets the color for all defined LEDs. UpdateColors needs to be called to make the changes visible
func (led *LED) SetColorAll(r byte, g byte, b byte) {
	for _, part := range led.GetParts() {
		led.SetColorPart(part, r, g, b)
	}
}

// Update makes color changes visible
func (led *LED) Update() error {
	if led.apa102Dummy == true && led.drawDummy == false {
		return nil
	}

	return led.apa102.Draw(led.apa102.Bounds(), led.image, image.Point{})
}

// getTotalNumLeds returns the number of leds (max LED ID + 1)
func (led *LED) getTotalNumLeds() int {
	return led.maxLedID + 1
}

// mapLedPartPos return the real LED ID based on the part name and position inside the part
func (led *LED) mapLedPartPos(part string, pos int) int {
	ledIDs, exists := led.partLedMap[part]
	if !exists {
		panic("invalid part name")
	}
	return ledIDs[pos]
}

// getRange returns an int array from `first` to `last`. It also works if `first` is bigger than `last`.
func getRange(first int, last int) []int {
	var length, direction int
	if first <= last {
		// forwards
		length = last - first + 1
		direction = +1
	} else {
		// backwards
		length = first - last + 1
		direction = -1
	}
	result := make([]int, length)
	for i := 0; i < length; i++ {
		result[i] = first + (direction * i)
	}
	return result
}
