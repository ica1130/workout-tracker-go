package main

import (
	"net/http"
	"strings"

	"workout-tracker-go.ilijakrilovic.com/internal/data"
)

var allowedCategories = map[string]bool{
	"shoulders": true,
	"chest":     true,
	"back":      true,
	"arms":      true,
	"core":      true,
	"legs":      true,
	"cardio":    true,
}

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

func (app *application) getExercisesByCategoryHandler(w http.ResponseWriter, r *http.Request) {

	category := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("category")))

	if category == "" || !allowedCategories[category] {
		http.Error(w, "invalid category", http.StatusBadRequest)
		return
	}

	exercises, err := app.models.Exercises.GetByCategory(category)
	if err != nil {
		http.Error(w, "the server encountered a problem and could not process your request", http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"exercises": exercises}, nil)
	if err != nil {
		http.Error(w, "the server encountered a problem and could not process your request", http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
	}
}

func (app *application) updateExerciseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	exercise, err := app.models.Exercises.GetById(id)
	if err != nil {
		switch {
		case err == data.ErrRecordNotFound:
			http.NotFound(w, r)
		default:
			http.Error(w, "the server encountered a problem and could not process your request", http.StatusInternalServerError)
			app.logger.Printf("error: %v", err)
		}
		return
	}

	var input struct {
		Name        string `json:"name"`
		Category    string `json:"category"`
		Description string `json:"description"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	exercise.Name = input.Name
	exercise.Category = input.Category
	exercise.Description = input.Description

	err = app.models.Exercises.Update(exercise)
	if err != nil {
		http.Error(w, "the server encountered a problem and could not process your request", http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"exercise": exercise}, nil)
	if err != nil {
		http.Error(w, "the server encountered a problem and could not process your request", http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
		return
	}
}
