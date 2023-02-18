package main

import (
	"encoding/json"
	"errors"
	"io"
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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	err := json.NewDecoder(r.Body).Decode(data)
	// TODO: triage those errors returned from Decode
	if err != nil {
		switch {
		case errors.Is(err, io.EOF):
			return errors.New("received empty body")
		default:
			return errors.New("received malformed json")
		}
	}
	return nil
}
