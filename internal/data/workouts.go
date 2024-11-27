package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Workout struct {
	ID       int64            `json:"id"`
	MemberID int64            `json:"member_id"`
	Date     time.Time        `json:"date"`
	Details  []*WorkoutDetail `json:"details"`
	Version  int              `json:"-"`
}

type WorkoutDetail struct {
	ID          int64   `json:"-"`
	WorkoutID   int64   `json:"workout_id"`
	ExerciseID  int64   `json:"exercise_id"`
	Set         int     `json:"set"`
	Repetitions int     `json:"repetitions"`
	Weight      float64 `json:"weight"`
}

type WorkoutModel struct {
	DB *sql.DB
}

func (w WorkoutModel) Insert(workout *Workout) error {

	tx, err := w.DB.Begin()
	if err != nil {
		return err
	}

	workoutQuery := `
		INSERT INTO workouts (member_id, date)
		VALUES ($1, $2)
		RETURNING id
	`

	args := []interface{}{workout.MemberID, workout.Date}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = tx.QueryRowContext(ctx, workoutQuery, args...).Scan(&workout.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	values := []string{}
	args = []interface{}{}
	argCounter := 1

	for _, detail := range workout.Details {
		detail.WorkoutID = workout.ID
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", argCounter, argCounter+1, argCounter+2, argCounter+3, argCounter+4))
		args = append(args, detail.WorkoutID, detail.ExerciseID, detail.Set, detail.Repetitions, detail.Weight)
		argCounter += 5
	}

	fmt.Println(values)

	detailsQuery := `
		INSERT INTO workout_details (workout_id, exercise_id, set, repetitions, weight)
		VALUES ` + strings.Join(values, ", ")

	_, err = tx.ExecContext(ctx, detailsQuery, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
