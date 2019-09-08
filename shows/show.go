package shows

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"

	"github.com/google/uuid"
	"github.com/light-bull/lightbull/hardware"
	"github.com/spf13/viper"
)

// Show is a collection of visuals
type Show struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	visuals       []*Visual
	currentVisual *Visual

	mux sync.Mutex
}

// showJSON is the format for a serialized JSON configuration
type showJSON struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	Visuals []*Visual `json:"visuals"`
}

// newShow creates a new show with the given name. It is meant to be called from ShowCollection.
func newShow(name string) (*Show, error) {
	if name == "" {
		return nil, errors.New("Invalid show name")
	}

	show := Show{ID: uuid.New(), Name: name} // FIXME: uuid is randomly generated, so there could be a collission

	return &show, nil
}

// MarshalJSON is there to implement the `json.Marshaller` interface.
func (show *Show) MarshalJSON() ([]byte, error) {
	data := showJSON{ID: show.ID, Name: show.Name, Visuals: show.visuals}
	return json.Marshal(data)
	//return json.MarshalIndent(data, "", "    ")
}

// UnmarshalJSON is there to implement the `json.Unmarshaller` interface.
func (show *Show) UnmarshalJSON(data []byte) error {
	input := showJSON{}

	err := json.Unmarshal(data, &input)
	if err != nil {
		return err
	}

	show.ID = input.ID
	show.Name = input.Name
	show.visuals = input.Visuals

	// TODO: input validation

	return nil
}

// newShowFromFile creates a new show and loads the data from the specified file. It is meant to be called from ShowCollection.
// FIXME: move somewhere else
func newShowFromFile(filepath string) (*Show, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	show := Show{}
	err = show.UnmarshalJSON(data)
	if err != nil {
		return nil, errors.New("Malformed show configuration")
	}

	return &show, nil
}

// Save writes the show configuration to disk
// FIXME: move somewhere else
func (show *Show) Save() error {
	show.mux.Lock()
	defer show.mux.Unlock()

	//data, err := show.MarshalJSON()
	data, err := json.MarshalIndent(show, "", "    ")
	if err != nil {
		log.Print("Error while serializing JSON for show")
		return err
	}

	err = ioutil.WriteFile(show.diskFile(), data, 0644)
	if err != nil {
		log.Print("Failed to write show to the config file: " + err.Error())
		return err
	}

	return nil
}

// delete deletes the show configuration from disk
// FIXME: move somewhere else
func (show *Show) delete() {
	show.mux.Lock()
	defer show.mux.Unlock()

	os.Remove(show.diskFile())
}

func (show *Show) diskFile() string {
	return path.Join(viper.GetString("directories.config"), "shows", show.ID.String()+".json")
}

// Visuals returns a list of all visuals
func (show *Show) Visuals() []*Visual {
	return show.visuals
}

// NewVisual creates a new visual in the show.
func (show *Show) NewVisual(name string) *Visual {
	visual := newVisual(name)

	show.mux.Lock()
	show.visuals = append(show.visuals, visual)
	show.mux.Unlock()

	show.Save() // TODO: asynchronly and throttled save show

	return visual
}

// DeleteVisual deletes the visual from the show.
func (show *Show) DeleteVisual(visual *Visual) {
	show.mux.Lock()
	defer show.mux.Unlock()

	for pos, cur := range show.visuals {
		if visual.ID == cur.ID {
			show.visuals = append(show.visuals[:pos], show.visuals[pos+1:]...)
			break
		}
	}
}

// CurrentVisual returns the visual that is currently played
func (show *Show) CurrentVisual() *Visual {
	return show.currentVisual
}

// SetCurrentVisual sets the visual that is currently played
func (show *Show) SetCurrentVisual(visual *Visual) {
	show.mux.Lock()
	defer show.mux.Unlock()

	show.currentVisual = visual
}

// Update decides about the changes that are caused by the current visual for a certain timestep.
func (show *Show) Update(hw *hardware.Hardware, nanoseconds int64) {
	visual := show.CurrentVisual()
	if visual != nil {
		visual.Update(hw, nanoseconds)
	}
}
