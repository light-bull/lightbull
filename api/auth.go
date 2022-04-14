package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/light-bull/lightbull/api/utils"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func (api *API) initAuth(router *mux.Router) {
	router.HandleFunc("/api/auth", api.handleAuth)
}

func (api *API) handleAuth(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)

	if r.Method == "POST" {
		// get data from request
		type inFormat struct {
			Password string `json:"password"`
		}
		data := inFormat{}
		err := utils.ParseJSON(&w, r, &data)
		if err != nil {
			return
		}

		// check password
		hash := viper.GetString("api.authentication")
		if bcrypt.CompareHashAndPassword([]byte(hash), []byte(data.Password)) == nil {
			jwt, err := api.jwt.New()
			if err != nil {
				log.Println("Cannot create JWT: " + err.Error())
				utils.WriteError(&w, "Internal server error", http.StatusInternalServerError)
			}

			type outFormat struct {
				Jwt string `json:"jwt"`
			}
			result := outFormat{Jwt: jwt}

			utils.WriteJSON(&w, result)
		} else {
			utils.WriteError(&w, "Invalid password", http.StatusUnauthorized)
		}
	} else {
		utils.WriteMethodNotAllowed(&w)
	}
}

func writeUnauthorized(w *http.ResponseWriter) {
	utils.WriteError(w, "Unauthorized", http.StatusUnauthorized)
}

// authenticate checks if a valid JWT is included in the requests and sends back an unauthorized error otherwise
func (api *API) authenticate(w *http.ResponseWriter, r *http.Request) bool {
	// get Authorization header
	authHeaders, ok := r.Header["Authorization"]
	if !ok {
		writeUnauthorized(w)
		return false
	}

	// we need exactly one Authorization header
	if len(authHeaders) != 1 {
		writeUnauthorized(w)
		return false
	}

	// get plain JWT (header is: "Bearer <jwt>")
	parts := strings.SplitN(authHeaders[0], " ", 2)
	if len(parts) != 2 {
		writeUnauthorized(w)
		return false
	}

	if parts[0] != "Bearer" {
		writeUnauthorized(w)
		return false
	}

	tokenString := parts[1]

	// check jwt
	if !api.jwt.Check(tokenString) {
		writeUnauthorized(w)
		return false
	}

	return true
}
