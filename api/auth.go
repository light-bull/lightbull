package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func (api *API) initAuth(router *mux.Router) {
	router.HandleFunc("/api/auth", api.handleAuth)
}

func (api *API) handleAuth(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "POST" {
		// get data from request
		type inFormat struct {
			Password string `json:"password"`
		}
		data := inFormat{}
		err := parseJSON(&w, r, &data)
		if err != nil {
			return
		}

		// check password
		hash := viper.GetString("api.authentication")
		if bcrypt.CompareHashAndPassword([]byte(hash), []byte(data.Password)) == nil {
			jwt, err := api.newJwt()
			if err != nil {
				log.Println("Cannot create JWT: " + err.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			type outFormat struct {
				Jwt string `json:"jwt"`
			}
			result := outFormat{Jwt: jwt}

			writeJSON(&w, result)
		} else {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
