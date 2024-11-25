package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/members", app.getMemberByEmailHandler)
	router.HandlerFunc(http.MethodPost, "/v1/members", app.createMemberHandler)
	router.HandlerFunc(http.MethodPut, "/v1/members/:id", app.updateMemberHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/members/:id", app.deleteMemberHandler)

	router.HandlerFunc(http.MethodPost, "/v1/exercises", app.createExerciseHandler)

	return router
}
