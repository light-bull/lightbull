package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/light-bull/lightbull/api/utils"
	"github.com/light-bull/lightbull/hardware"
)

func (api *API) initSystem(router *mux.Router) {
	router.HandleFunc("/api/shutdown", api.handleShutdown)
	router.HandleFunc("/api/ethernet", api.handleEthernet)
}

func (api *API) handleShutdown(w http.ResponseWriter, r *http.Request) {
	if !api.authenticate(&w, r) {
		return
	}
	utils.EnableCors(&w)

	if r.Method == "POST" {
		api.hw.System.Shutdown()
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (api *API) handleEthernet(w http.ResponseWriter, r *http.Request) {
	if !api.authenticate(&w, r) {
		return
	}
	utils.EnableCors(&w)

	if r.Method == "GET" {
		data := api.hw.System.EthernetConfig()

		utils.WriteJSON(&w, data)
	} else if r.Method == "PUT" {
		data := hardware.EthernetConfig{}

		err := utils.ParseJSON(&w, r, &data)
		if err != nil {
			return
		}

		err = api.hw.System.SetEthernetConfig(data)
		if err != nil {
			http.Error(w, "Invalid data: "+err.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
