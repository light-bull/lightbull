package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/light-bull/lightbull/api/mapper"
	"github.com/light-bull/lightbull/api/utils"
	"github.com/light-bull/lightbull/events"
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
	router.HandleFunc("/api/parameters/{id}/links", api.handleParameterLinks)
	router.HandleFunc("/api/parameters/{id}/links/{linkedId}", api.handleParameterLinks)

	router.HandleFunc("/api/current", api.handleCurrent)
}

func (api *API) handleShows(w http.ResponseWriter, r *http.Request) {
	if !api.authenticate(&w, r) {
		return
	}
	utils.EnableCors(&w)

	if r.Method == "GET" {
		utils.WriteJSON(&w, mapper.MapShows(api.shows.Shows()))
	} else if r.Method == "POST" {
		// get data from request
		type format struct {
			Name     string `json:"name"`
			Favorite bool   `json:"favorite"`
		}
		data := format{}
		err := utils.ParseJSON(&w, r, &data)
		if err != nil {
			return
		}

		// create new show
		show, err := api.shows.NewShow(data.Name, data.Favorite)
		if err != nil {
			utils.WriteError(&w, err.Error(), http.StatusBadRequest)
			return
		}

		show.Favorite = data.Favorite

		api.eventhub.PublishNew(events.ShowAdded, show, show, utils.GetConnectionID(r))

		// return show data, especially the ID may be interesting
		utils.WriteJSONWithStatus(&w, mapper.MapShowWithVisualIds(show), http.StatusCreated)
	} else {
		utils.WriteMethodNotAllowed(&w)
	}
}

func (api *API) handleShowDetails(w http.ResponseWriter, r *http.Request) {
	if !api.authenticate(&w, r) {
		return
	}
	utils.EnableCors(&w)

	// get show
	vars := mux.Vars(r)
	id := vars["id"]

	show := api.shows.FindShow(id)
	if show == nil {
		utils.WriteError(&w, "Invalid or unknown ID", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		utils.WriteJSON(&w, mapper.MapShowWithVisuals(show))
	} else if r.Method == "PUT" {
		// get data from request
		type format struct {
			Name     string `json:"name"`
			Favorite bool   `json:"favorite"`
		}
		data := format{}
		err := utils.ParseJSON(&w, r, &data)
		if err != nil {
			return
		}

		if data.Name != "" {
			show.Name = data.Name
		}

		show.Favorite = data.Favorite

		api.eventhub.PublishNew(events.ShowChanged, show, show, utils.GetConnectionID(r))
		utils.WriteJSON(&w, mapper.MapShowWithVisuals(show))
	} else if r.Method == "DELETE" {
		api.shows.DeleteShow(show)
		api.eventhub.PublishNew(events.ShowDeleted, show, show, utils.GetConnectionID(r))
		w.WriteHeader(http.StatusNoContent)
	} else {
		utils.WriteMethodNotAllowed(&w)
	}
}

func (api *API) handleVisuals(w http.ResponseWriter, r *http.Request) {
	if !api.authenticate(&w, r) {
		return
	}
	utils.EnableCors(&w)

	if r.Method == "GET" {
		utils.WriteJSON(&w, mapper.MapAllVisualsFromShows(api.shows.Shows()))
	} else if r.Method == "POST" {
		// get data from request
		type format struct {
			Name   string `json:"name"`
			ShowId string `json:"showId"`
		}
		data := format{}
		err := utils.ParseJSON(&w, r, &data)
		if err != nil {
			return
		}

		// get show
		show := api.shows.FindShow(data.ShowId)
		if show == nil {
			utils.WriteError(&w, "Invalid or unknown show ID", http.StatusBadRequest)
			return
		}

		// add visual to show
		visual := show.NewVisual(data.Name)
		api.eventhub.PublishNew(events.VisualAdded, visual, show, utils.GetConnectionID(r))

		utils.WriteJSONWithStatus(&w, mapper.MapVisualWithGroups(show.ID, visual), http.StatusCreated)
	} else {
		utils.WriteMethodNotAllowed(&w)
	}
}

func (api *API) handleVisualDetails(w http.ResponseWriter, r *http.Request) {
	if !api.authenticate(&w, r) {
		return
	}
	utils.EnableCors(&w)

	// get visual and show
	vars := mux.Vars(r)
	id := vars["id"]

	show, visual := api.shows.FindVisual(id)

	if visual == nil {
		utils.WriteError(&w, "Invalid or unknown ID", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		utils.WriteJSON(&w, mapper.MapVisualWithGroups(show.ID, visual))
	} else if r.Method == "PUT" {
		// get data from request
		type format struct {
			Name string `json:"name"`
		}
		data := format{}
		err := utils.ParseJSON(&w, r, &data)
		if err != nil {
			return
		}

		if data.Name != "" {
			visual.Name = data.Name
		}

		api.eventhub.PublishNew(events.VisualChanged, visual, show, utils.GetConnectionID(r))

		utils.WriteJSON(&w, mapper.MapVisualWithGroups(show.ID, visual))
	} else if r.Method == "DELETE" {
		show.DeleteVisual(visual)
		api.eventhub.PublishNew(events.VisualDeleted, visual, show, utils.GetConnectionID(r))
		w.WriteHeader(http.StatusNoContent)
	} else {
		utils.WriteMethodNotAllowed(&w)
	}
}

func (api *API) handleGroups(w http.ResponseWriter, r *http.Request) {
	if !api.authenticate(&w, r) {
		return
	}
	utils.EnableCors(&w)

	if r.Method == "POST" {
		// get data from request
		type format struct {
			VisualId   string   `json:"visualId"`
			Parts      []string `json:"parts"`
			EffectType string   `json:"effectType"`
		}
		data := format{}
		err := utils.ParseJSON(&w, r, &data)
		if err != nil {
			return
		}

		// get visual
		show, visual := api.shows.FindVisual(data.VisualId)
		if visual == nil {
			utils.WriteError(&w, "Invalid or unknown visual ID", http.StatusBadRequest)
			return
		}

		// add group
		group, err := visual.NewGroup(data.Parts, data.EffectType)
		if err != nil {
			utils.WriteError(&w, "Failed to create group: "+err.Error(), http.StatusBadRequest)
			return
		}

		api.eventhub.PublishNew(events.GroupAdded, group, show, utils.GetConnectionID(r))

		utils.WriteJSON(&w, mapper.MapGroup(visual.ID, group))
	} else {
		utils.WriteMethodNotAllowed(&w)
	}
}

func (api *API) handleGroupDetails(w http.ResponseWriter, r *http.Request) {
	if !api.authenticate(&w, r) {
		return
	}
	utils.EnableCors(&w)

	// get visual and show
	vars := mux.Vars(r)
	id := vars["id"]

	show, visual, group := api.shows.FindGroup(id)
	if group == nil {
		utils.WriteError(&w, "Invalid or unknown ID", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		utils.WriteJSON(&w, mapper.MapGroup(visual.ID, group))
	} else if r.Method == "PUT" {
		// get data from request
		type format struct {
			Parts      []string `json:"parts"`
			EffectType string   `json:"effectType"`
		}
		data := format{}
		err := utils.ParseJSON(&w, r, &data)
		if err != nil {
			return
		}

		if len(data.Parts) != 0 {
			err = group.SetParts(data.Parts)
			if err != nil {
				utils.WriteError(&w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		if data.EffectType != "" {
			err = group.SetEffect(data.EffectType)
			if err != nil {
				utils.WriteError(&w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		api.eventhub.PublishNew(events.GroupChanged, group, show, utils.GetConnectionID(r))
		utils.WriteJSON(&w, mapper.MapGroup(visual.ID, group))
	} else if r.Method == "DELETE" {
		visual.DeleteGroup(group)
		api.eventhub.PublishNew(events.GroupDeleted, group, show, utils.GetConnectionID(r))
		w.WriteHeader(http.StatusNoContent)
	} else {
		utils.WriteMethodNotAllowed(&w)
	}
}

func (api *API) handleParameterDetails(w http.ResponseWriter, r *http.Request) {
	if !api.authenticate(&w, r) {
		return
	}
	utils.EnableCors(&w)

	// get parameter
	vars := mux.Vars(r)
	id := vars["id"]

	show, _, _, parameter := api.shows.FindParameter(id)
	if parameter == nil {
		utils.WriteError(&w, "Invalid or unknown ID", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		utils.WriteJSON(&w, parameter)
	} else if r.Method == "PUT" {
		// get data from request
		type format struct {
			Current *json.RawMessage `json:"current"`
			Default *json.RawMessage `json:"default"`
		}
		data := format{}
		err := utils.ParseJSON(&w, r, &data)
		if err != nil {
			return
		}

		// change current value (if given)
		eventTopic := ""
		if data.Current != nil {
			err = parameter.SetFromJSON(*data.Current)
			if err != nil {
				utils.WriteError(&w, "Failed to set parameter: "+err.Error(), http.StatusBadRequest)
				return
			}
			eventTopic = events.ParameterChanged
		}

		// change default value (if given)
		if data.Default != nil {
			err = parameter.SetDefaultFromJSON(*data.Default)
			if err != nil {
				utils.WriteError(&w, "Failed to set parameter: "+err.Error(), http.StatusBadRequest)
				return
			}

			eventTopic = events.ParameterDefaultChanged
		}

		api.eventhub.PublishNew(eventTopic, parameter, show, utils.GetConnectionID(r))
	} else {
		utils.WriteMethodNotAllowed(&w)
	}
}

func (api *API) handleParameterLinks(w http.ResponseWriter, r *http.Request) {
	if !api.authenticate(&w, r) {
		return
	}
	utils.EnableCors(&w)

	// get first parameter
	vars := mux.Vars(r)
	id := vars["id"]
	linkedId, linkedIdFound := vars["linkedId"]

	show, visual1, _, parameter1 := api.shows.FindParameter(id)
	if parameter1 == nil {
		utils.WriteError(&w, "Invalid or unknown ID", http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		// we do not accept the linked id in the url here, this is only allowed for DELETE
		if linkedIdFound {
			w.WriteHeader(http.StatusNotFound)
		}

		// get data from request
		type format struct {
			LinkedParameter uuid.UUID `json:"linkedParameter"`
		}
		data := format{}
		err := utils.ParseJSON(&w, r, &data)
		if err != nil {
			return
		}

		// get second parameter
		_, visual2, _, parameter2 := api.shows.FindParameter(data.LinkedParameter.String())
		if parameter2 == nil {
			utils.WriteError(&w, "Invalid or unknown ID for linked parameter", http.StatusBadRequest)
			return
		}

		// validate that parameters belong to same visual
		if visual1.ID != visual2.ID {
			utils.WriteError(&w, "Parameters must belong to same visual", http.StatusBadRequest)
			return
		}

		err = visual1.LinkParameter(parameter1, parameter2)
		if err != nil {
			utils.WriteError(&w, err.Error(), http.StatusBadRequest)
			return
		}

		// publish changes on both parameters
		api.eventhub.PublishNew(events.ParameterLinksChanged, parameter1, show, utils.GetConnectionID(r))
		api.eventhub.PublishNew(events.ParameterLinksChanged, parameter2, show, utils.GetConnectionID(r))
		w.WriteHeader(http.StatusOK)
	} else if r.Method == "DELETE" {
		// we need the linkedId from the URL
		if !linkedIdFound {
			w.WriteHeader(http.StatusNotFound)
		}

		// get second parameter
		_, _, _, parameter2 := api.shows.FindParameter(linkedId)
		if parameter2 == nil {
			utils.WriteError(&w, "Invalid or unknown ID for linked parameter", http.StatusBadRequest)
			return
		}

		err := visual1.UnlinkParameter(parameter1, parameter2)
		if err != nil {
			utils.WriteError(&w, err.Error(), http.StatusBadRequest)
			return
		}

		// publish changes on both parameters
		api.eventhub.PublishNew(events.ParameterLinksChanged, parameter1, show, utils.GetConnectionID(r))
		api.eventhub.PublishNew(events.ParameterLinksChanged, parameter2, show, utils.GetConnectionID(r))
		w.WriteHeader(http.StatusOK)
	} else {
		utils.WriteMethodNotAllowed(&w)
	}
}

func (api *API) handleCurrent(w http.ResponseWriter, r *http.Request) {
	if !api.authenticate(&w, r) {
		return
	}
	utils.EnableCors(&w)

	if r.Method == "GET" {
		utils.WriteJSON(&w, api.getCurrentShowAndVisual())
	} else if r.Method == "PUT" {
		// get data
		type format struct {
			ShowId   string `json:"showId"`
			VisualId string `json:"visualId"`
		}
		data := format{}
		err := utils.ParseJSON(&w, r, &data)
		if err != nil {
			return
		}

		// get show and visual
		var show *shows.Show
		var visual *shows.Visual

		if data.ShowId != "" {
			show = api.shows.FindShow(data.ShowId)
			if show == nil {
				utils.WriteError(&w, "Invalid or unknown show ID", http.StatusNotFound)
				return
			}
		}

		if data.VisualId != "" {
			_, visual = api.shows.FindVisual(data.VisualId)
			if visual == nil {
				utils.WriteError(&w, "Invalid or unknown visual ID", http.StatusNotFound)
				return
			}
		}

		// set current show and visual
		err = api.shows.SetCurrentVisual(show, visual)
		if err != nil {
			utils.WriteError(&w, err.Error(), http.StatusBadRequest)
			return
		}

		currentShowAndVisual := api.getCurrentShowAndVisual()

		// Send event
		api.eventhub.PublishNew(events.CurrentChanged, currentShowAndVisual, nil, utils.GetConnectionID(r))
		utils.WriteJSON(&w, currentShowAndVisual)
	} else if r.Method == "DELETE" {
		if api.shows.CurrentShow() != nil {
			api.shows.ClearCurrentVisual()
		}

		utils.WriteJSON(&w, api.getCurrentShowAndVisual())
	} else {
		utils.WriteMethodNotAllowed(&w)
	}
}

// getCurrentShowAndVisual returns a struct with the current show and visual
func (api *API) getCurrentShowAndVisual() mapper.CurrentShowAndVisualJSON {
	show, visual := api.shows.GetCurrentVisual()
	return mapper.MapCurrent(show, visual)
}
