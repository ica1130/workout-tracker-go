package main

import (
	"net/http"
	"time"

	"workout-tracker-go.ilijakrilovic.com/internal/data"
)

func (app *application) createWorkoutHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		MemberID int64                 `json:"member_id"`
		Date     time.Time             `json:"date"`
		Details  []*data.WorkoutDetail `json:"details"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	workout := &data.Workout{
		MemberID: input.MemberID,
		Date:     input.Date,
		Details:  input.Details,
	}

	err = app.models.Workouts.Insert(workout)
	if err != nil {
		http.Error(w, "there was an error while creating an exercise", http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"workout": workout}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.logger.Printf("error: %v", err)
		return
	}
}
