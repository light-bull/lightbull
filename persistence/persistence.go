package persistence

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/light-bull/lightbull/events"
	"github.com/light-bull/lightbull/shows"
	"github.com/spf13/viper"
)

// Persistence stores dynamic configuration and shows on disk
type Persistence struct {
	configDir string
	showsDir  string

	eventhub    *events.EventHub
	eventclient *EventClient
}

// NewPersistence returns a new persistence store
func NewPersistence(eventhub *events.EventHub) (*Persistence, error) {
	persistence := Persistence{}

	persistence.configDir = viper.GetString("directories.config")
	if err := os.MkdirAll(persistence.configDir, 0755); err != nil {
		return nil, errors.New("Cannot create configuration directory: " + err.Error())
	}

	persistence.showsDir = path.Join(persistence.configDir, "shows")
	if err := os.MkdirAll(persistence.showsDir, 0755); err != nil {
		return nil, errors.New("Cannot create configuration directory for shows: " + err.Error())
	}

	persistence.eventhub = eventhub
	persistence.eventclient = newEventClient(&persistence)
	persistence.eventhub.RegisterClient(persistence.eventclient)

	return &persistence, nil
}

// SaveConfig stores the serialized configuration as JSON on disk
func (persistence *Persistence) SaveConfig(name string, data interface{}, secret bool) error {
	file := path.Join(persistence.configDir, name+".json")
	return persistence.save(file, data, secret)
}

// LoadConfig loads a JSON config file and deserializes it
func (persistence *Persistence) LoadConfig(name string, data interface{}) error {
	file := path.Join(persistence.configDir, name+".json")
	return persistence.load(file, data)
}

// HasConfig checks if a configuration file exists. It does not look into the file.
func (persistence *Persistence) HasConfig(name string) bool {
	file := path.Join(persistence.configDir, name+".json")
	if info, err := os.Stat(file); err == nil && !info.IsDir() {
		return true
	}
	return false
}

// SaveShow stores the given show on disk
func (persistence *Persistence) SaveShow(show *shows.Show) error {
	// TODO: mutex!
	file := path.Join(persistence.showsDir, show.ID.String()+".json")
	return persistence.save(file, show, false)
}

// DeleteShow deletes the show from disk
func (persistence *Persistence) DeleteShow(show *shows.Show) {
	file := path.Join(persistence.showsDir, show.ID.String()+".json")
	os.Remove(file)
}

// LoadShows loads all shows from disk and adds them to the show collection
func (persistence *Persistence) LoadShows(showCollection *shows.ShowCollection) {
	files, _ := filepath.Glob(persistence.showsDir + "/*.json")
	if files == nil {
		log.Print("No shows loaded.")
		return
	}

	for _, path := range files {
		show := shows.Show{}
		err := persistence.load(path, &show)
		if err != nil {
			log.Print("Error while loading show from file: " + path + "(" + err.Error() + ")")
			continue
		}

		showCollection.AppendShow(&show)
	}

	return
}

// save stores the serialized object as JSON on disk
func (persistence *Persistence) save(file string, data interface{}, secret bool) error {
	// serialize
	jsonOutput, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	// secret files should only be readable for own user
	var mode os.FileMode = 0644
	if secret {
		mode = 0600
	}

	// write
	err = ioutil.WriteFile(file, jsonOutput, mode)
	if err != nil {
		return err
	}

	return nil
}

// load loads a JSON file and deserializes it
func (persistence *Persistence) load(file string, data interface{}) error {
	jsonInput, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonInput, &data)
	if err != nil {
		return err
	}

	return nil
}
