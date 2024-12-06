package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Members   MemberModel
	Exercises ExerciseModel
	Workouts  WorkoutModel
	Tokens    TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Members:   MemberModel{DB: db},
		Exercises: ExerciseModel{DB: db},
		Workouts:  WorkoutModel{DB: db},
		Tokens:    TokenModel{DB: db},
	}
}
