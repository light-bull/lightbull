package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/light-bull/lightbull/shows/effects"
)

func (api *API) initConfig(router *mux.Router) {
	router.HandleFunc("/api/config", api.handleConfig)
}

func (api *API) handleConfig(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "GET" {
		type format struct {
			Parts   []string          `json:"parts"`
			Effects map[string]string `json:"effects"`
		}
		data := format{
			Parts:   api.hw.Led.GetParts(),
			Effects: effects.GetEffects(),
		}

		writeJSON(&w, data)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
