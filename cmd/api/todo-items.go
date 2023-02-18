package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/kickbu2towski/unicorn-be/cmd/api/internal/data"
	"github.com/kickbu2towski/unicorn-be/cmd/api/internal/validator"
)

func (app *application) showTodoItem(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	tI := data.TodoItem{
		Id:       id,
		Name:     "Read 100 pages of Let's Go Further book",
		State:    "NEXT",
		Priority: "B",
		Tags:     []string{"personal"},
	}

	err = app.writeJSON(w, envelope{"todoItem": tI}, nil, http.StatusOK)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createTodoItem(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string   `json:"name"`
		State    string   `json:"state"`
		Priority string   `json:"priority"`
		Tags     []string `json:"tags"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err.Error())
		return
	}

	todoItem := data.TodoItem{
		Name:     input.Name,
		State:    input.State,
		Priority: input.Priority,
		Tags:     input.Tags,
	}

	v := validator.New()
	data.ValidateTodoItem(v, &todoItem)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}
