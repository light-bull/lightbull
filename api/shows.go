package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/light-bull/lightbull/shows"
)

func (api *API) initShows(router *mux.Router) {
	router.HandleFunc("/api/shows", api.handleShows)
	router.HandleFunc("/api/shows/{id}", api.handleShowDetails)

	router.HandleFunc("/api/visuals", api.handleVisuals)
	router.HandleFunc("/api/visuals/{id}", api.handleVisualDetails)

	router.HandleFunc("/api/groups", api.handleGroups)
	router.HandleFunc("/api/groups/{id}", api.handleGroupDetails)

	router.HandleFunc("/api/parameters/{id}", api.handleParameterDetails)

	router.HandleFunc("/api/current", api.handleCurrent)
}

func (api *API) handleShows(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "GET" {
		shows := api.shows.Shows()

		type format struct {
			ID       uuid.UUID   `json:"id"`
			Name     string      `json:"name"`
			Favorite bool        `json:"favorite"`
			Visuals  []uuid.UUID `json:"visuals"`
		}
		data := make([]format, len(shows))

		for i, show := range shows {
			data[i].ID = show.ID
			data[i].Name = show.Name
			data[i].Favorite = show.Favorite
			data[i].Visuals = make([]uuid.UUID, len(show.Visuals()))
			for j, visual := range show.Visuals() {
				data[i].Visuals[j] = visual.ID
			}
		}

		writeJSON(&w, data)
	} else if r.Method == "POST" {
		// get data from request
		type format struct {
			Name     string `json:"name"`
			Favorite bool   `json:"favorite"`
		}
		data := format{}
		err := parseJSON(&w, r, &data)
		if err != nil {
			return
		}

		// create new show
		show, err := api.shows.NewShow(data.Name, data.Favorite)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		show.Favorite = data.Favorite
		show.Save()

		// return show data, especially the ID may be interesting
		writeJSON(&w, show)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (api *API) handleShowDetails(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// get show
	vars := mux.Vars(r)
	id := vars["id"]

	show := api.shows.FindShow(id)
	if show == nil {
		http.Error(w, "Invalid or unknown ID", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		type visualFormat struct {
			ID   uuid.UUID `json:"id"`
			Name string    `json:"name"`
		}

		type showFormat struct {
			ID       uuid.UUID      `json:"id"`
			Name     string         `json:"name"`
			Favorite bool           `json:"favorite"`
			Visuals  []visualFormat `json:"visuals"`
		}

		data := showFormat{
			ID:       show.ID,
			Name:     show.Name,
			Favorite: show.Favorite,
			Visuals:  make([]visualFormat, len(show.Visuals())),
		}
		for i, visual := range show.Visuals() {
			data.Visuals[i].ID = visual.ID
			data.Visuals[i].Name = visual.Name
		}

		writeJSON(&w, data)
	} else if r.Method == "PUT" {
		// get data from request
		type format struct {
			Name     string `json:"name"`
			Favorite bool   `json:"favorite"`
		}
		data := format{}
		err := parseJSON(&w, r, &data)
		if err != nil {
			return
		}

		if data.Name != "" {
			show.Name = data.Name
		}

		show.Favorite = data.Favorite

		// TODO: move (async) save to shows.Show
		show.Save()
	} else if r.Method == "DELETE" {
		api.shows.DeleteShow(show)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (api *API) handleVisuals(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "GET" {
		type visualFormat struct {
			Show uuid.UUID `json:"show"`
			ID   uuid.UUID `json:"id"`
			Name string    `json:"name"`
		}

		var result []visualFormat
		for _, show := range api.shows.Shows() {
			for _, visual := range show.Visuals() {
				visual := visualFormat{
					Show: show.ID,
					ID:   visual.ID,
					Name: visual.Name,
				}
				result = append(result, visual)
			}
		}
		writeJSON(&w, result)
	} else if r.Method == "POST" {
		// get data from request
		type format struct {
			Name string `json:"name"`
			Show string `json:"show"`
		}
		data := format{}
		err := parseJSON(&w, r, &data)
		if err != nil {
			return
		}

		// get show
		show := api.shows.FindShow(data.Show)
		if show == nil {
			http.Error(w, "Invalid or unknown show ID", http.StatusBadRequest)
			return
		}

		// add visual to show
		show.NewVisual(data.Name)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (api *API) handleVisualDetails(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// get visual and show
	vars := mux.Vars(r)
	id := vars["id"]

	show, visual := api.shows.FindVisual(id)

	if visual == nil {
		http.Error(w, "Invalid or unknown ID", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		writeJSON(&w, visual)
	} else if r.Method == "PUT" {
		// get data from request
		type format struct {
			Name string `json:"name"`
		}
		data := format{}
		err := parseJSON(&w, r, &data)
		if err != nil {
			return
		}

		if data.Name != "" {
			visual.Name = data.Name
		}

		// TODO: move (async) save to shows.Show
		show.Save()
	} else if r.Method == "DELETE" {
		show.DeleteVisual(visual)
		show.Save() // TODO: do this somewhere else
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (api *API) handleGroups(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "POST" {
		// get data from request
		type format struct {
			Visual string   `json:"visual"`
			Parts  []string `json:"parts"`
			Effect string   `json:"effect"`
		}
		data := format{}
		err := parseJSON(&w, r, &data)
		if err != nil {
			return
		}

		// get visual
		show, visual := api.shows.FindVisual(data.Visual)
		if visual == nil {
			http.Error(w, "Invalid or unknown visual ID", http.StatusBadRequest)
			return
		}

		// add group
		_, err = visual.NewGroup(data.Parts, data.Effect)
		if err != nil {
			http.Error(w, "Failed to create group: "+err.Error(), http.StatusBadRequest)
		}

		// TODO: move (async) save to shows.Show
		show.Save()
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (api *API) handleGroupDetails(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// get visual and show
	vars := mux.Vars(r)
	id := vars["id"]

	show, visual, group := api.shows.FindGroup(id)
	if group == nil {
		http.Error(w, "Invalid or unknown ID", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		writeJSON(&w, group)
	} else if r.Method == "PUT" {
		// get data from request
		type format struct {
			Parts  []string `json:"parts"`
			Effect string   `json:"effect"`
		}
		data := format{}
		err := parseJSON(&w, r, &data)
		if err != nil {
			return
		}

		if len(data.Parts) != 0 {
			group.SetParts(data.Parts)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		}

		if data.Effect != "" {
			group.SetEffect(data.Effect)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		}

		// TODO: move (async) save to shows.Show
		show.Save()
	} else if r.Method == "DELETE" {
		visual.DeleteGroup(group)
		show.Save() // TODO: do this somewhere else
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (api *API) handleParameterDetails(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// get parameter
	vars := mux.Vars(r)
	id := vars["id"]

	show, _, _, parameter := api.shows.FindParameter(id)
	if parameter == nil {
		http.Error(w, "Invalid or unknown ID", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		writeJSON(&w, parameter)
	} else if r.Method == "PUT" {
		// get data from request
		type format struct {
			Current *json.RawMessage `json:"current"`
			Default *json.RawMessage `json:"default"`
		}
		data := format{}
		err := parseJSON(&w, r, &data)
		if err != nil {
			return
		}

		if data.Current != nil {
			err = parameter.SetFromJSON(*data.Current)
			if err != nil {
				http.Error(w, "Failed to set parameter: "+err.Error(), http.StatusBadRequest)
			}
		}

		if data.Default != nil {
			err = parameter.SetDefaultFromJSON(*data.Default)
			if err != nil {
				http.Error(w, "Failed to set parameter: "+err.Error(), http.StatusBadRequest)
			}

			// TODO: move (async) save to shows.Show
			show.Save()
		}

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (api *API) handleCurrent(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "GET" {
		type format struct {
			Show   *uuid.UUID `json:"show"`
			Visual *uuid.UUID `json:"visual"`
		}
		data := format{}

		show := api.shows.CurrentShow()
		if show != nil {
			data.Show = &show.ID

			visual := show.CurrentVisual()
			if visual != nil {
				data.Visual = &visual.ID
			}
		}

		writeJSON(&w, data)
	} else if r.Method == "PUT" {
		// get data
		type format struct {
			Show   string `json:"show"`
			Visual string `json:"visual"`
			Blank  bool   `json:"blank"`
		}
		data := format{}
		err := parseJSON(&w, r, &data)
		if err != nil {
			return
		}

		// handle blank first because its easy
		if data.Blank == true {
			if api.shows.CurrentShow() != nil {
				api.shows.CurrentShow().SetCurrentVisual(nil)
			}

			return
		}

		// get show and visual
		var show, showOfVisual *shows.Show
		var visual *shows.Visual

		if data.Show != "" {
			show = api.shows.FindShow(data.Show)
			if show == nil {
				http.Error(w, "Invalid or unknown show ID", http.StatusBadRequest)
				return
			}
		}

		if data.Visual != "" {
			showOfVisual, visual = api.shows.FindVisual(data.Visual)
			if visual == nil {
				http.Error(w, "Invalid or unknown visual ID", http.StatusBadRequest)
				return
			}
		}

		// set current show and visual
		if show != nil && visual != nil {
			// show and visual given -> check that visual belongs to show
			if show != showOfVisual {
				http.Error(w, "Visual does not belong to show", http.StatusBadRequest)
				return
			}

			show.SetCurrentVisual(visual)
			api.shows.SetCurrentShow(show)
		} else if show != nil && visual == nil {
			// only show given -> set show and no current visual
			if show != api.shows.CurrentShow() {
				show.SetCurrentVisual(nil)
				api.shows.SetCurrentShow(show)
			}
		} else if show == nil && visual != nil {
			// only visual -> check that visual belongs to current show
			if showOfVisual != api.shows.CurrentShow() {
				http.Error(w, "Visual does not belong to current show", http.StatusBadRequest)
				return
			}

			api.shows.CurrentShow().SetCurrentVisual(visual)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
