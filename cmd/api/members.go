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
