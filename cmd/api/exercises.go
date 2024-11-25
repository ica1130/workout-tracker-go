package main

import (
	"net/http"

	"workout-tracker-go.ilijakrilovic.com/internal/data"
)

func (app *application) createExerciseHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name        string `json:"name"`
		Category    string `json:"category"`
		Description string `json:"description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.logger.Printf("error: %v", err)
		return
	}

	exercise := &data.Exercise{
		Name:        input.Name,
		Category:    input.Category,
		Description: input.Description,
	}

	err = app.models.Exercises.Insert(exercise)
	if err != nil {
		http.Error(w, "there was an error while creating an exercise", http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"exercise": exercise}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
		return
	}
}
