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
		http.Error(w, "email must be provided", http.StatusBadRequest)
		return
	}

	member, err := app.models.Members.GetByEmail(email)
	if err != nil {
		http.Error(w, "error while retreiving member", http.StatusInternalServerError)
		app.logger.Printf("error while retreiving member: %v", err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"members": member}, nil)
	if err != nil {
		http.Error(w, "error while retreiving member", http.StatusInternalServerError)
		app.logger.Printf("error while printing member: %v", err)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.logger.Printf("error: %v", err)
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
		http.Error(w, "error: server encountered an error while processing your request", http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
		return
	}

	err = app.models.Members.Insert(member)
	if err != nil {
		http.Error(w, "there was an error while creating a member", http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
		return
	}

	token, err := app.models.Tokens.New(member.ID, 3*24*time.Hour, "activation")
	if err != nil {
		http.Error(w, "error: server encountered an error while processing your request", http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
		return
	}

	responseEnvelope := envelope{
		"member":           member,
		"activation_token": token.Plaintext,
	}

	err = app.writeJSON(w, http.StatusCreated, responseEnvelope, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
		return
	}
}

func (app *application) updateMemberHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	member, err := app.models.Members.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			http.NotFound(w, r)
			app.logger.Printf("error: %v", err)
		default:
			http.Error(w, "error: server encountered an error while processing your request", http.StatusInternalServerError)
			app.logger.Printf("error: %v", err)
		}
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
		http.Error(w, "error: bad reuqest", http.StatusBadRequest)
		return
	}

	member.Email = input.Email
	member.Name = input.Name
	member.Height = input.Height
	member.Weight = input.Weight

	err = app.models.Members.Update(member)
	if err != nil {
		http.Error(w, "error: server encountered an error while processing your request", http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"member": member}, nil)
	if err != nil {
		http.Error(w, "error: server encountered an error while processing your request", http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
		return
	}
}

func (app *application) deleteMemberHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	err = app.models.Members.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			http.NotFound(w, r)
		default:
			http.Error(w, "error: server encountered an error while processing your request", http.StatusInternalServerError)
			app.logger.Printf("error: %v", err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "member sucessfully deleted"}, nil)
	if err != nil {
		http.Error(w, "error: server encountered an error while processing your request", http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
		return
	}
}
