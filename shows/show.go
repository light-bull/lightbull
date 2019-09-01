package shows

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"path"
	"sync"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// Show is a collection of visuals
type Show struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	visuals []*Visual

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

	data, err := show.MarshalJSON()
	if err != nil {
		log.Print("Error while serializing JSON for show")
		return err
	}

	file := path.Join(viper.GetString("directories.config"), "shows", show.ID.String()+".json")
	err = ioutil.WriteFile(file, data, 0644)
	if err != nil {
		log.Print("Failed to write show to the config file: " + err.Error())
		return err
	}

	return nil
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
