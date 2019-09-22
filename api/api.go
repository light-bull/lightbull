package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"

	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/shows"
	_ "github.com/light-bull/lightbull/statik" // needed for statik
)

// API implements the REST API
type API struct {
	hw    *hardware.Hardware
	shows *shows.ShowCollection
}

// New starts the listener for the REST API
func New(hw *hardware.Hardware, shows *shows.ShowCollection) (*API, error) {
	api := API{hw: hw, shows: shows}
	router := mux.NewRouter()

	// API
	api.initConfig(router)
	api.initSystem(router)
	api.initShows(router)

	// Frontend
	statikFS, err := fs.New()
	if err != nil {
		return nil, err
	}

	router.PathPrefix("/web").Handler(http.StripPrefix("/web/", http.FileServer(statikFS)))

	// Redirect to frontend
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/", http.StatusTemporaryRedirect)
	})

	// Run server
	port := viper.GetInt("listen")
	if port <= 0 || port > 65535 {
		return nil, errors.New("Invalid port to listen on")
	}
	go http.ListenAndServe(fmt.Sprintf(":%d", port), router)

	return &api, nil
}
