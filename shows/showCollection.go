package shows

import (
	"log"
	"path"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
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
func (showCollection *ShowCollection) NewShow(name string) (*Show, error) {
	showCollection.mux.Lock()
	defer showCollection.mux.Unlock()

	show, err := newShow(name)
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

// SetCurrentShow set the show that is currently played
func (showCollection *ShowCollection) SetCurrentShow(show *Show) {
	showCollection.mux.Lock()
	defer showCollection.mux.Unlock()

	showCollection.currentShow = show
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
