package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/kickbu2towski/unicorn-be/cmd/api/internal/data"
	"github.com/kickbu2towski/unicorn-be/cmd/api/internal/validator"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Authentication")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		values := strings.Split(authHeader, " ")
		if len(values) != 2 || values[0] != "Bearer" {
			app.invalidAuthTokenResponse(w, r)
			return
		}

		token := values[1]

		v := validator.New()
		if data.ValidateToken(v, token); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		user, err := app.models.UserModel.GetForToken(token, data.ScopeAuthentication)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrNoRecord):
				app.invalidAuthTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user.IsAnonymousUser() {
			app.failedAuthResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (app *application) requireUserActivated(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}

	return app.requireAuth(fn)
}
