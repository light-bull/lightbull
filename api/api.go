package api

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"

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

	jwtKey []byte
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
	err := api.initJwt()
	if err != nil {
		return nil, err
	}

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

// initJwt prepares the key material to issue JWTs
func (api *API) initJwt() error {
	type format struct {
		Key       []byte    `json:"key"`
		Generated time.Time `json:"generated"`
	}
	data := format{}

	if api.persistence.HasConfig("jwt") {
		// config is there -> load it of fail
		if err := api.persistence.Load("jwt", &data); err != nil {
			return err
		}

		api.jwtKey = data.Key
	} else {
		// generate key and store it
		api.jwtKey = make([]byte, 64)
		if _, err := rand.Read(api.jwtKey); err != nil {
			return errors.New("Failed to generate secret key for JWT: " + err.Error())
		}

		data.Key = api.jwtKey
		data.Generated = time.Now()
		if err := api.persistence.Save("jwt", &data, true); err != nil {
			return err
		}
	}

	return nil
}

// newJwt issues a new JSON Web Token
func (api *API) newJwt() (string, error) {
	claims := &jwt.StandardClaims{
		Issuer:    jwtIssuer,
		ExpiresAt: time.Now().Add(jwtValidity).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(api.jwtKey)
}

// checkJwt validates the given JSON Web Token
func (api *API) checkJwt(tokenString string) bool {
	// parse the token
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}

		// give signing key to parser
		return api.jwtKey, nil
	})

	if !token.Valid {
		return false
	}

	// get claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	// validate claims
	if claims.Valid() != nil {
		return false
	}

	return true
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
	if ok := api.checkJwt(tokenString); !ok {
		http.Error(*w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	return true
}
