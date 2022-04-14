package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/viper"

	"github.com/gorilla/mux"

	"github.com/light-bull/lightbull/api/utils"
	"github.com/light-bull/lightbull/events"
	"github.com/light-bull/lightbull/frontend"
	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/persistence"
	"github.com/light-bull/lightbull/shows"
)

// API implements the REST API
type API struct {
	hw          *hardware.Hardware
	shows       *shows.ShowCollection
	eventhub    *events.EventHub
	persistence *persistence.Persistence
	jwt         *utils.JWTManager
}

// New starts the listener for the REST API
func New(hw *hardware.Hardware, shows *shows.ShowCollection, eventhub *events.EventHub, persistence *persistence.Persistence) (*API, error) {
	api := API{
		hw:          hw,
		shows:       shows,
		eventhub:    eventhub,
		persistence: persistence,
	}

	router := mux.NewRouter()

	// Get key material for JWTs
	jwtManager, err := utils.NewJWTManager(persistence)
	if err != nil {
		return nil, err
	}
	api.jwt = jwtManager

	// API routes
	api.initAuth(router)
	api.initConfig(router)
	api.initSystem(router)
	api.initShows(router)
	api.initWS(router)

	// Frontend
	router.PathPrefix("/web").Handler(http.StripPrefix("/web/", http.FileServer(http.FS(frontend.Frontend))))

	// Redirect to frontend
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/", http.StatusTemporaryRedirect)
	})

	// Run server
	port := viper.GetInt("api.listen")
	if port <= 0 || port > 65535 {
		return nil, errors.New("Invalid port to listen on")
	}
	go http.ListenAndServe(fmt.Sprintf(":%d", port), router)

	return &api, nil
}
