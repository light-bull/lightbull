package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/light-bull/lightbull/shows/effects"
)

func (api *API) initShows(router *mux.Router) {
	router.HandleFunc("/api/shows", api.handleShows)
	router.HandleFunc("/api/shows/{id}", api.handleShowDetails)

	router.HandleFunc("/api/visuals", api.handleVisuals)
	router.HandleFunc("/api/visuals/{id}", api.handleVisualDetails)

	router.HandleFunc("/api/groups/{id}", api.handleGroup)
}

func (api *API) handleShows(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "GET" {
		shows := api.shows.Shows()

		type format struct {
			ID       uuid.UUID `json:"id"`
			Name     string    `json:"name"`
			Favorite bool      `json:"favorite"`
		}
		data := make([]format, len(shows))

		for i, show := range shows {
			data[i].ID = show.ID
			data[i].Name = show.Name
			data[i].Favorite = show.Favorite
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

	if r.Method == "POST" {
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
		visual := show.NewVisual(data.Name)

		// testing things
		group := visual.NewGroup()
		group.SetParts(api.hw.Led.GetParts())
		group.SetEffect(effects.SingleColor)
		api.shows.SetCurrentShow(show)
		show.SetCurrentVisual(visual)
		show.Save()
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
	} else if r.Method == "POST" {
		// TODO
	} else if r.Method == "DELETE" {
		show.DeleteVisual(visual)
		show.Save() // TODO: do this somewhere else
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (api *API) handleGroup(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// get show
	vars := mux.Vars(r)
	id := vars["id"]

	_, _, group := api.shows.FindGroup(id)

	if group == nil {
		http.Error(w, "Invalid or unknown ID", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		writeJSON(&w, group)
	} else if r.Method == "POST" {
		// TODO
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
