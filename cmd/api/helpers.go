package main

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, data any, headers http.Header, status int) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')

	for key, val := range headers {
		w.Header()[key] = val
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)
	w.Write(js)
	return nil
}
