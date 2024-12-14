package main

import (
	"errors"
	"net/http"
	"strings"

	"workout-tracker-go.ilijakrilovic.com/internal/data"
	"workout-tracker-go.ilijakrilovic.com/internal/validator"
)

func (app *application) createExerciseHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name        string `json:"name"`
		Category    string `json:"category"`
		Description string `json:"description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	exercise := &data.Exercise{
		Name:        input.Name,
		Category:    input.Category,
		Description: input.Description,
	}

	v := validator.New()

	if data.ValidateExercise(v, exercise); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Exercises.Insert(exercise)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"exercise": exercise}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getExercisesByCategoryHandler(w http.ResponseWriter, r *http.Request) {

	category := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("category")))

	v := validator.New()

	if data.ValidateCategory(v, category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	exercises, err := app.models.Exercises.GetByCategory(category)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"exercises": exercises}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateExerciseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	exercise, err := app.models.Exercises.GetById(id)
	if err != nil {
		switch {
		case err == data.ErrRecordNotFound:
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
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
		app.badRequestResponse(w, r, err)
		return
	}

	exercise.Name = input.Name
	exercise.Category = input.Category
	exercise.Description = input.Description

	err = app.models.Exercises.Update(exercise)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"exercise": exercise}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteExerciseHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
	}

	err = app.models.Exercises.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"exercise": "exercise sucessfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
