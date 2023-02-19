package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) Routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/checkhealth", app.checkhealth)
	router.HandlerFunc(http.MethodPost, "/v1/todoItems", app.createTodoItem)
	router.HandlerFunc(http.MethodGet, "/v1/todoItems/:id", app.showTodoItem)
	router.HandlerFunc(http.MethodDelete, "/v1/todoItems/:id", app.deleteTodoItem)
	router.HandlerFunc(http.MethodPut, "/v1/todoItems/:id", app.updateTodoItem)

	return router
}
