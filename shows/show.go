package shows

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/light-bull/lightbull/hardware"
)

// Show is a collection of visuals
type Show struct {
	ID       uuid.UUID
	Name     string
	Favorite bool

	visuals       []*Visual
	currentVisual *Visual

	mux sync.Mutex
}

// showJSON is the format for a serialized JSON configuration
type showJSON struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Favorite bool      `json:"favorite"`

	Visuals []*Visual `json:"visuals"`
}

// newShow creates a new show with the given name. It is meant to be called from ShowCollection.
func newShow(name string, favorite bool) (*Show, error) {
	if name == "" {
		return nil, errors.New("Invalid show name")
	}

	show := Show{ID: uuid.New(), Name: name, Favorite: favorite} // FIXME: uuid is randomly generated, so there could be a collission

	return &show, nil
}

// MarshalJSON is there to implement the `json.Marshaller` interface.
func (show *Show) MarshalJSON() ([]byte, error) {
	data := showJSON{ID: show.ID, Name: show.Name, Favorite: show.Favorite, Visuals: show.visuals}
	return json.Marshal(data)
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
	show.Favorite = input.Favorite
	show.visuals = input.Visuals

	// TODO: input validation

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

// setCurrentVisual sets the visual that is currently played
func (show *Show) setCurrentVisual(visual *Visual) {
	show.mux.Lock()
	defer show.mux.Unlock()

	show.currentVisual = visual
}

// hasVisual checks if the visual belongs to the show
func (show *Show) hasVisual(visual *Visual) bool {
	show.mux.Lock()
	defer show.mux.Unlock()

	for _, cur := range show.visuals {
		if cur == visual {
			return true
		}
	}
	return false
}

// Update decides about the changes that are caused by the current visual for a certain timestep.
func (show *Show) Update(hw *hardware.Hardware, nanoseconds int64) {
	visual := show.CurrentVisual()
	if visual != nil {
		visual.Update(hw, nanoseconds)
	} else {
		hw.Led.SetColorAll(0, 0, 0)
	}
}
