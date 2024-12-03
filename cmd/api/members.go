package main

import (
	"errors"
	"net/http"
	"time"

	"workout-tracker-go.ilijakrilovic.com/internal/data"
)

func (app *application) getMemberByEmailHandler(w http.ResponseWriter, r *http.Request) {

	email := r.URL.Query().Get("email")

	if email == "" {
		app.badRequestResponse(w, r, errors.New("missing email parameter"))
		return
	}

	member, err := app.models.Members.GetByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"members": member}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) createMemberHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Height   int64  `json:"height"`
		Weight   int64  `json:"weight"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	member := &data.Member{
		Email:     input.Email,
		Name:      input.Name,
		Activated: false,
		Height:    input.Height,
		Weight:    input.Weight,
	}

	err = member.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Members.Insert(member)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(member.ID, 3*24*time.Hour, "activation")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	responseEnvelope := envelope{
		"member":           member,
		"activation_token": token.Plaintext,
	}

	err = app.writeJSON(w, http.StatusCreated, responseEnvelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateMemberHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Email  string `json:"email"`
		Name   string `json:"name"`
		Height int64  `json:"height"`
		Weight int64  `json:"weight"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	member, err := app.models.Members.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	member.Email = input.Email
	member.Name = input.Name
	member.Height = input.Height
	member.Weight = input.Weight

	err = app.models.Members.Update(member)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"member": member}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteMemberHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
	}

	err = app.models.Members.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "member sucessfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) activateMemberHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlain string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	member, err := app.models.Members.GetForToken("activation", input.TokenPlain)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	member.Activated = true

	err = app.models.Members.Update(member)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Tokens.DeleteAllForMember("activation", member.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"member": member}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
