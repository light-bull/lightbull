package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/light-bull/lightbull/api/utils"
)

func (api *API) initSimulator(router *mux.Router) {
	router.HandleFunc("/api/simulator", api.handleSimulator)
}

func (api *API) handleSimulator(w http.ResponseWriter, r *http.Request) {
	if !api.authenticate(&w, r) {
		return
	}
	utils.EnableCors(&w)

	if r.Method == "GET" {
		type colorformat struct {
			R byte `json:"r"`
			G byte `json:"g"`
			B byte `json:"b"`
		}

		data := make(map[string][]colorformat)
		for _, part := range api.hw.Led.GetParts() {
			data[part] = make([]colorformat, api.hw.Led.GetNumLeds((part)))

			for i := 0; i < api.hw.Led.GetNumLeds(part); i++ {
				r, g, b := api.hw.Led.GetColor(part, i)
				data[part][i].R = r
				data[part][i].G = g
				data[part][i].B = b
			}
		}

		utils.WriteJSON(&w, data)
	} else {
		utils.WriteMethodNotAllowed(&w)
	}
}
