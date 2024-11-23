package main

import (
	"net/http"
	"time"

	"workout-tracker-go.ilijakrilovic.com/internal/data"
)

func (app *application) getMembersHandler(w http.ResponseWriter, r *http.Request) {
	member := &data.Member{
		ID:        1,
		Email:     "ilijakrilovic@gmail.com",
		Name:      "Ilija",
		Height:    192,
		Weight:    95,
		CreatedAt: time.Now(),
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"members": member}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
