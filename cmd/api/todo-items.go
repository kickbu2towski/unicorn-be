package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kickbu2towski/unicorn-be/cmd/api/internal/data"
	"github.com/kickbu2towski/unicorn-be/cmd/api/internal/validator"
)

func (app *application) showTodoItem(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	todoItem, err := app.models.TodoItems.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, envelope{"todoItem": todoItem}, nil, http.StatusOK)
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

	err = app.models.TodoItems.Insert(&todoItem)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", todoItem.Id))

	err = app.writeJSON(w, envelope{"todoItem": todoItem}, headers, http.StatusCreated)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTodoItem(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.TodoItems.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, envelope{"message": "todo item delete successfully"}, nil, http.StatusOK)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTodoItem(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	todoItem, err := app.models.TodoItems.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name     string     `json:"name"`
		State    string     `json:"state"`
		Priority string     `json:"priority"`
		Tags     []string   `json:"tags"`
		ClosedAt *time.Time `json:"closed_at"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err.Error())
		return
	}

	todoItem.Name = input.Name
	todoItem.State = input.State
	todoItem.Priority = input.Priority
	todoItem.Tags = input.Tags
	todoItem.ClosedAt = input.ClosedAt

	v := validator.New()
	data.ValidateTodoItem(v, todoItem)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.TodoItems.Update(id, todoItem)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, envelope{"todoItem": todoItem}, nil, http.StatusOK)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
