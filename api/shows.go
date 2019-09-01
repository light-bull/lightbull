package api

import (
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

	router.HandleFunc("/api/groups/{id}", api.handleGroup)
}

func (api *API) handleShows(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	type format struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

	if r.Method == "GET" {
		shows := api.shows.Shows()

		data := make([]format, len(shows))

		for i, show := range shows {
			data[i].ID = show.ID
			data[i].Name = show.Name
		}

		writeJSON(&w, data)
	} else if r.Method == "POST" || r.Method == "PUT" {
		// get data from request
		data := format{}
		err := parseJSON(&w, r, &data)
		if err != nil {
			return
		}

		// create or update show
		if r.Method == "POST" {
			// create new show
			show, err := api.shows.NewShow(data.Name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			writeJSON(&w, show)
		} else if r.Method == "PUT" {
			// update existing show
			show := api.shows.FindShow(data.ID.String())
			if show == nil {
				http.Error(w, "Invalid or unknown ID", http.StatusBadRequest)
				return
			}

			// TODO: move validation and (async) save to shows.Show
			if data.Name == "" {
				http.Error(w, "Invalid name", http.StatusBadRequest)
				return
			}

			show.Name = data.Name
			show.Save()
		}
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
		/*
			// TODO
			visual := show.NewVisual("testvisual")
			visual.NewGroup()
			visual.NewGroup()
		*/

		// FIXME: only return show and visual data, no groups/effects/...

		writeJSON(&w, show)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (api *API) handleVisuals(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "POST" || r.Method == "PUT" {
		// FIXME

		// get data from request
		data := shows.Visual{}
		err := parseJSON(&w, r, &data)
		if err != nil {
			return
		}
		/*
			if r.Method == "POST" || r.Method == "PUT" {
					// get data from request
					data := shows.Show{}
					err := parseJSON(&w, r, &data)
					if err != nil {
						return
					}

					// create or update show
					if r.Method == "POST" {
						// create new show
						show, err := api.shows.NewShow(data.Name)
						if err != nil {
							http.Error(w, err.Error(), http.StatusBadRequest)
							return
						}

						writeJSON(&w, show)
					} else if r.Method == "PUT" {
						// update existing show
						show := api.shows.FindShow(data.ID.String())
						if show == nil {
							http.Error(w, "Invalid or unknown ID", http.StatusBadRequest)
							return
						}

						err := show.Update(&data)
						if err != nil {
							http.Error(w, err.Error(), http.StatusBadRequest)
							return
						}
					}
				} else {
					http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
				}
		*/
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (api *API) handleVisualDetails(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// get visual and show
	vars := mux.Vars(r)
	id := vars["id"]

	_, visual := api.shows.FindVisual(id)

	if visual == nil {
		http.Error(w, "Invalid or unknown ID", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		writeJSON(&w, visual)
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
