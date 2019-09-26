package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// WriteJSON serializes the given struct to JSON and sends it back over the HTTP connection
func WriteJSON(w *http.ResponseWriter, v interface{}) {
	result, err := json.Marshal(v)
	if err != nil {
		http.Error(*w, "Error while serializing JSON", http.StatusInternalServerError)
	}

	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(result)
}

// ParseJSON reads the body from the HTTP request and tries to deserialize it as JSON
func ParseJSON(w *http.ResponseWriter, r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(*w, "Error reading request body", http.StatusInternalServerError)
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		http.Error(*w, "Invalid data format", http.StatusBadRequest)
		return err
	}

	return nil
}