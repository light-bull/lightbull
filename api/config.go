package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/light-bull/lightbull/api/utils"
	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows/effects"
)

func (api *API) initConfig(router *mux.Router) {
	router.HandleFunc("/api/config", api.handleConfig)
}

func (api *API) handleConfig(w http.ResponseWriter, r *http.Request) {
	if !api.authenticate(&w, r) {
		return
	}
	utils.EnableCors(&w)

	if r.Method == "GET" {
		type format struct {
			Parts    []string          `json:"parts"`
			Effects  map[string]string `json:"effects"`
			Features []string          `json:"features"`
		}

		data := format{
			Parts:   api.hw.Led.GetParts(),
			Effects: effects.GetEffects(),
			Features: make([]string, 0),
		}

		if api.hw.System.EthernetConfig().Mode != hardware.EthUnmanaged {
			data.Features = append(data.Features, "ethernet")
		}

		utils.WriteJSON(&w, data)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
