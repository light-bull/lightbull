package utils

import "net/http"

// EnableCors sets the Access-Control-Allow-Origin header to "*"
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
