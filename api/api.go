package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"

	"github.com/light-bull/lightbull/api/utils"
	"github.com/light-bull/lightbull/events"
	"github.com/light-bull/lightbull/hardware"
	"github.com/light-bull/lightbull/persistence"
	"github.com/light-bull/lightbull/shows"
	_ "github.com/light-bull/lightbull/statik" // needed for statik
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
	api.initWebsocket(router)

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
	port := viper.GetInt("api.listen")
	if port <= 0 || port > 65535 {
		return nil, errors.New("Invalid port to listen on")
	}
	go http.ListenAndServe(fmt.Sprintf(":%d", port), router)

	return &api, nil
}

// authenticate checks if a valid JWT is included in the requests and sends back an unauthorized error otherwise
func (api *API) authenticate(w *http.ResponseWriter, r *http.Request) bool {
	// get Authorization header
	authHeaders, ok := r.Header["Authorization"]
	if !ok {
		http.Error(*w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	// we need exactly one Authorization header
	if len(authHeaders) != 1 {
		http.Error(*w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	// get plain JWT (header is: "Bearer <jwt>")
	parts := strings.SplitN(authHeaders[0], " ", 2)
	if len(parts) != 2 {
		http.Error(*w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	if parts[0] != "Bearer" {
		http.Error(*w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	tokenString := parts[1]

	// check jwt
	if !api.jwt.Check(tokenString) {
		http.Error(*w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	return true
}
