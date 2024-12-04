package main

import (
	"context"
	"net/http"

	"workout-tracker-go.ilijakrilovic.com/internal/data"
)

type contextKey string

const memberContextKey = contextKey("member")

func (app *application) contextSetMember(r *http.Request, member *data.Member) *http.Request {
	ctx := context.WithValue(r.Context(), memberContextKey, member)
	return r.WithContext(ctx)
}

func (app *application) contextGetMember(r *http.Request) *data.Member {
	member, ok := r.Context().Value(memberContextKey).(*data.Member)
	if !ok {
		panic("member not found in context")
	}

	return member
}
