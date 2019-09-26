package utils

import (
	"net/http"

	"github.com/google/uuid"
)

// EnableCors sets the Access-Control-Allow-Origin header to "*"
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// GetConnectionID gets the websocket connection id from the `X-Lightbull-Connection-Id` header
func GetConnectionID(r *http.Request) uuid.UUID {
	// get all X-Lightbull-Connection-Id header
	headers, ok := r.Header["X-Lightbull-Connection-Id"]
	if !ok {
		return uuid.Nil
	}

	// we need exactly one such header
	if len(headers) != 1 {
		return uuid.Nil
	}

	// check if it is a valid uuid
	id, err := uuid.Parse(headers[0])
	if err != nil {
		return uuid.Nil
	}

	return id
}
