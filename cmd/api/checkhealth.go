package main

import (
	"net/http"
)

func (app *application) checkhealth(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"version": version,
		"system_info": map[string]string{
			"status": "OK",
			"env":    app.config.env,
		},
	}

	err := app.writeJSON(w, data, nil, 200)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
