package main

import (
	"context"
	"net/http"

	"github.com/kickbu2towski/unicorn-be/cmd/api/internal/data"
)

type contextKey string

const userContext = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContext, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContext).(*data.User)
	if !ok {
		panic("missing user in request context")
	}
	return user
}
