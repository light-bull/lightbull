package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func writeJSON(w *http.ResponseWriter, v interface{}) {
	result, err := json.Marshal(v)
	if err != nil {
		http.Error(*w, "Error while serializing JSON", http.StatusInternalServerError)
	}

	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(result)
}

func parseJSON(w *http.ResponseWriter, r *http.Request, v interface{}) error {
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
