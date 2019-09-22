package shows

import (
	"errors"
	"log"
	"path"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/light-bull/lightbull/shows/parameters"
	"github.com/spf13/viper"
)

// ShowCollection is the collection of all available shows
type ShowCollection struct {
	shows       []*Show
	currentShow *Show

	mux sync.Mutex
}

// NewShowCollection loads all shows from the configuration file and returns the new ShowCollection objects
func NewShowCollection() *ShowCollection {
	showCollection := ShowCollection{}

	showCollection.loadShows()

	return &showCollection
}

// Shows returns a list of loaded shows
func (showCollection *ShowCollection) Shows() []*Show {
	return showCollection.shows
}

// NewShow creates a new show
func (showCollection *ShowCollection) NewShow(name string, favorite bool) (*Show, error) {
	showCollection.mux.Lock()
	defer showCollection.mux.Unlock()

	show, err := newShow(name, favorite)
	if err != nil {
		return nil, err
	}

	showCollection.shows = append(showCollection.shows, show)

	err = show.Save()
	if err != nil {
		return nil, err
	}

	return show, nil
}

// DeleteShow creates a new show
func (showCollection *ShowCollection) DeleteShow(show *Show) {
	showCollection.mux.Lock()
	defer showCollection.mux.Unlock()

	// delete from disk
	// TODO: trigger delete based on events
	show.delete()

	// delete from list
	for pos, cur := range showCollection.shows {
		if show.ID == cur.ID {
			showCollection.shows = append(showCollection.shows[:pos], showCollection.shows[pos+1:]...)
			break
		}
	}
}

// CurrentShow returns the show that is currently played
func (showCollection *ShowCollection) CurrentShow() *Show {
	return showCollection.currentShow
}

// GetCurrentVisual returns the current visual and show
func (showCollection *ShowCollection) GetCurrentVisual() (show *Show, visual *Visual) {
	if showCollection.currentShow != nil {
		return showCollection.currentShow, showCollection.currentShow.CurrentVisual()
	}
	return nil, nil
}

// SetCurrentVisual set the show and visual that is currently played.
// If the show is `nil`, the visual needs to belong to the current show.
// If the show is changed and no visual is given, the current visual is always set to `nil`.
func (showCollection *ShowCollection) SetCurrentVisual(show *Show, visual *Visual) error {
	showCollection.mux.Lock()
	defer showCollection.mux.Unlock()

	if show != nil && visual != nil {
		// show and visual given -> check that visual belongs to show
		if show.hasVisual(visual) == false {
			return errors.New("Visual does not belong to show")
		}

		show.setCurrentVisual(visual)
		showCollection.currentShow = show
	} else if show != nil && visual == nil {
		// only show given -> set show and no current visual
		if showCollection.currentShow != show {
			show.setCurrentVisual(nil)
			showCollection.currentShow = show
		}
	} else if show == nil && visual != nil {
		// only visual -> check that visual belongs to current show (and that there is a current show)
		if showCollection.currentShow == nil {
			return errors.New("No current show set, cannot set a new visual")
		} else if showCollection.currentShow.hasVisual(visual) == false {
			return errors.New("Visual does not belong to the current show")
		}

		showCollection.currentShow.setCurrentVisual(visual)
	} else {
		return errors.New("Visual or show need to be specified")
	}
	return nil
}

// ClearCurrentVisual resets the current visual to `nil` so that nothing is played
func (showCollection *ShowCollection) ClearCurrentVisual() {
	showCollection.mux.Lock()
	defer showCollection.mux.Unlock()

	if showCollection.currentShow != nil {
		showCollection.currentShow.setCurrentVisual(nil)
	}
}

// FindShow returns the show with the given ID or nil for malformed and non-existing IDs
func (showCollection *ShowCollection) FindShow(idStr string) *Show {
	// Parse UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil
	}

	// Locking
	showCollection.mux.Lock()
	defer showCollection.mux.Unlock()

	// search for show
	for _, show := range showCollection.shows {
		if show.ID == id {
			return show
		}
	}

	return nil
}

// FindVisual returns the visual with the given ID and the beloning show or nil for malformed and non-existing IDs
func (showCollection *ShowCollection) FindVisual(idStr string) (*Show, *Visual) {
	// Parse UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, nil
	}

	// Locking
	showCollection.mux.Lock()
	defer showCollection.mux.Unlock()

	// iterate over shows and visuals
	for _, show := range showCollection.shows {
		for _, visual := range show.Visuals() {
			if visual.ID == id {
				return show, visual
			}
		}
	}

	return nil, nil
}

// FindGroup returns the group with the given ID and the beloning show and visual or nil for malformed and non-existing IDs
func (showCollection *ShowCollection) FindGroup(idStr string) (*Show, *Visual, *Group) {
	// Parse UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, nil, nil
	}

	// Locking
	showCollection.mux.Lock()
	defer showCollection.mux.Unlock()

	// iterate over shows and visuals
	for _, show := range showCollection.shows {
		for _, visual := range show.Visuals() {
			for _, group := range visual.Groups() {
				if group.ID == id {
					return show, visual, group
				}
			}
		}
	}

	return nil, nil, nil
}

// FindParameter returns the parameter with the given ID and the belonging show, visual and group or nil for malformed and non-existing IDs
func (showCollection *ShowCollection) FindParameter(idStr string) (*Show, *Visual, *Group, *parameters.Parameter) {
	// Parse UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, nil, nil, nil
	}

	// Locking
	showCollection.mux.Lock()
	defer showCollection.mux.Unlock()

	// iterate over shows, visuals and groups
	for _, show := range showCollection.shows {
		for _, visual := range show.Visuals() {
			for _, group := range visual.Groups() {
				for _, parameter := range group.Effect.Parameters() {
					if parameter.ID == id {
						return show, visual, group, parameter
					}
				}
			}
		}
	}

	return nil, nil, nil, nil
}

// loadShows loads the stored shows from the configuration files
func (showCollection *ShowCollection) loadShows() {
	showCollection.mux.Lock()
	defer showCollection.mux.Unlock()

	dir := path.Join(viper.GetString("directories.config"), "shows")
	files, _ := filepath.Glob(dir + "/*.json")
	if files == nil {
		log.Print("No shows loaded.")
		return
	}

	for _, path := range files {
		show, err := newShowFromFile(path)
		if err != nil {
			log.Print("Error while loading show from file: " + path + "(" + err.Error() + ")")
			continue
		}

		showCollection.shows = append(showCollection.shows, show)
	}
}
