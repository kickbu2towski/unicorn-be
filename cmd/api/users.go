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
	user.GenerateHash()

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
			ActivationToken: token.PlainTextToken,
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
	var input struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err.Error())
		return
	}

	v := validator.New()
	if data.ValidateToken(v, input.Token); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.UserModel.GetForToken(input.Token, data.ScopeActivation)
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

	err = app.models.TokenModel.DeleteAllForUser(user.ID, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user.Activated = true
	err = app.models.UserModel.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, envelope{"message": "user activation successful"}, nil, http.StatusOK)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createAuthToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err.Error())
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePassword(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.UserModel.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			app.invalidAuthCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	matches := user.PasswordMatches(input.Password)
	if !matches {
		app.invalidAuthCredentialsResponse(w, r)
		return
	}

	expiry := 24 * time.Hour
	token, err := app.models.TokenModel.New(user.ID, expiry, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, envelope{"auth_token": token}, nil, http.StatusOK)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
