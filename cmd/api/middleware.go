package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
	"workout-tracker-go.ilijakrilovic.com/internal/data"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetMember(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if !claims.Valid(time.Now()) {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if claims.Issuer != "workout-tracker-go.ilijakrilovic.com" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if !claims.AcceptAudience("workout-tracker-go.ilijakrilovic.com") {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		memberID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		member, err := app.models.Members.GetById(memberID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetMember(r, member)

		next.ServeHTTP(w, r)
	})
}
