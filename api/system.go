package api

import (
	"net/http"

	"github.com/light-bull/lightbull/hardware"
	"github.com/gorilla/mux"
)

func (api *API) initSystem(router *mux.Router) {
	router.HandleFunc("/api/shutdown", api.handleShutdown)
	router.HandleFunc("/api/ethernet", api.handleEthernet)
}

func (api *API) handleShutdown(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "POST" {
		api.hw.System.Shutdown()
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (api *API) handleEthernet(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "GET" {
		data := api.hw.System.EthernetConfig()

		writeJSON(&w, data)
	} else if r.Method == "PUT" {
		data := hardware.EthernetConfig{}

		err := parseJSON(&w, r, &data)
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
