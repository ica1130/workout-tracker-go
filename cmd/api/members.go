package main

import (
	"net/http"

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
		Email  string `json:"email"`
		Name   string `json:"name"`
		Height int64  `json:"height"`
		Weight int64  `json:"weight"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.logger.Printf("error: %v", err)
		return
	}

	member := &data.Member{
		Email:  input.Email,
		Name:   input.Name,
		Height: input.Height,
		Weight: input.Weight,
	}

	err = app.models.Members.Insert(member)
	if err != nil {
		http.Error(w, "there was an error while creating a member", http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"member": member}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
	}
}
