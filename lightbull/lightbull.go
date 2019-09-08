package lightbull

import (
	"time"

	"github.com/light-bull/lightbull/api"
	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows"
	"github.com/spf13/viper"
)

// Lightbull contains all software components (`hardware`, `api`, `shows`).
// This basically is the glue code between the other packages.
type Lightbull struct {
	Hardware *hardware.Hardware
	Shows    *shows.ShowCollection
	API      *api.API
}

// New prepares the whole lightbull controller for use: it initializes the hardware, starts the
// hardware update loop and starts the REST API.
func New() (*Lightbull, error) {
	lightbull := Lightbull{}
	var err error

	// TODO: create directories that are needed

	// initialize hardware and run update loop
	lightbull.Hardware, err = hardware.New()
	if err != nil {
		return nil, err
	}

	// load shows
	lightbull.Shows = shows.NewShowCollection()

	// run update loop for modes and hardware
	go lightbull.UpdateLoop()

	// run api server
	lightbull.API, err = api.New(lightbull.Hardware, lightbull.Shows)
	if err != nil {
		return nil, err
	}

	return &lightbull, nil
}

// UpdateLoop runs the current mode program and writes changes to the hardware in regular intervals
func (lightbull *Lightbull) UpdateLoop() {
	lastUpdate := time.Now()
	sleepTime := time.Duration(1000000000.0 / viper.GetFloat64("leds.fps"))
	for {
		time.Sleep(sleepTime)

		nanoseconds := time.Since(lastUpdate).Nanoseconds()
		lastUpdate = time.Now()
		show := lightbull.Shows.CurrentShow()
		if show != nil {
			show.Update(lightbull.Hardware, nanoseconds)
		}

		lightbull.Hardware.Update()
	}
}
