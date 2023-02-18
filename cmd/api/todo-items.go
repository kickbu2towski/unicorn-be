package main

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/kickbu2towski/unicorn-be/cmd/api/internal/data"
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
	w.Write([]byte("creating a new todo item\n"))
}
