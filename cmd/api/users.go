package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/kickbu2towski/unicorn-be/cmd/api/internal/data"
	"github.com/kickbu2towski/unicorn-be/cmd/api/internal/validator"
)

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err.Error())
		return
	}

	user := &data.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.UserModel.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := app.models.TokenModel.New(user.ID, (3*24)*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		tmplData := struct {
			ActivationToken string
			ID              int64
		}{
			ActivationToken: token,
			ID:              user.ID,
		}
		err := app.mailer.Send(user.Email, "activate_user.tmpl", &tmplData)
		if err != nil {
			app.logger.Print(err)
		}
	})

	err = app.writeJSON(w, envelope{"user": user}, nil, http.StatusCreated)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUser(w http.ResponseWriter, r *http.Request) {
	// 1. read json
	var input struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err.Error())
		return
	}

	// 2. validate token
	v := validator.New()
	if data.ValidateToken(v, input.Token); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// 3. get user for token
	user, err := app.models.UserModel.GetForToken(input.Token)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			v.AddError("token", "invalid token or token expired")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// 4. activate user
	user.Activated = true
	err = app.models.UserModel.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// 5. send response
	err = app.writeJSON(w, envelope{"message": "user activation successful"}, nil, http.StatusOK)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
