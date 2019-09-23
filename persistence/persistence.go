package persistence

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/viper"
)

// Persistence stores dynamic configuration and shows on disk
type Persistence struct {
	configDir string
	showsDir  string
}

// NewPersistence returns a new persistence store
func NewPersistence() (*Persistence, error) {
	persistence := Persistence{}

	persistence.configDir = viper.GetString("directories.config")
	if err := os.MkdirAll(persistence.configDir, 0755); err != nil {
		return nil, errors.New("Cannot create configuration directory: " + err.Error())
	}

	persistence.showsDir = path.Join(persistence.configDir, "shows")
	if err := os.MkdirAll(persistence.showsDir, 0755); err != nil {
		return nil, errors.New("Cannot create configuration directory for shows: " + err.Error())
	}
	return &persistence, nil
}

// Save stores the serialized object as JSON on disk
func (persistence *Persistence) Save(name string, data interface{}, secret bool) error {
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
	file := path.Join(persistence.configDir, name+".json")
	err = ioutil.WriteFile(file, jsonOutput, mode)
	if err != nil {
		return err
	}

	return nil
}

// Load loads a JSON file and deserializes it
func (persistence *Persistence) Load(name string, data interface{}) error {
	file := path.Join(persistence.configDir, name+".json")
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

// HasConfig checks if a configuration file exists. It does not look into the file.
func (persistence *Persistence) HasConfig(name string) bool {
	file := path.Join(persistence.configDir, name+".json")
	if info, err := os.Stat(file); err == nil && !info.IsDir() {
		return true
	}
	return false
}
