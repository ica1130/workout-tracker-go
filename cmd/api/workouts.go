package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
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
		app.badRequestResponse(w, r, err)
		return
	}

	workout := &data.Workout{
		MemberID: input.MemberID,
		Date:     input.Date,
		Details:  input.Details,
	}

	err = app.models.Workouts.Insert(workout)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"workout": workout}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAllWorkoutsByMemberIDHandler(w http.ResponseWriter, r *http.Request) {
	memberID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	workouts, err := app.models.Workouts.GetByMemberID(memberID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"workouts": workouts}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteWorkoutHandler(w http.ResponseWriter, r *http.Request) {

	id := httprouter.ParamsFromContext(r.Context()).ByName("workout_id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Workouts.Delete(idInt)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "workout sucessfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
