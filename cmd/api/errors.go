package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Print(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, message any, status int) {
	data := envelope{
		"error": message,
	}

	err := app.writeJSON(w, data, nil, status)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, message any) {
	app.errorResponse(w, r, message, http.StatusBadRequest)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered an error and could not process the request"
	app.errorResponse(w, r, message, http.StatusInternalServerError)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, message, http.StatusNotFound)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method isn't allowed on this resource", r.Method)
	app.errorResponse(w, r, message, http.StatusMethodNotAllowed)
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, err map[string]string) {
	app.errorResponse(w, r, err, http.StatusUnprocessableEntity)
}

func (app *application) invalidAuthCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "received invalid auth credentials"
	app.errorResponse(w, r, message, http.StatusUnauthorized)
}
