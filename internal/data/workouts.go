package data

import (
	"context"
	"database/sql"
	"time"
)

type Workout struct {
	ID       int64     `json:"id"`
	MemberID int64     `json:"member_id"`
	Date     time.Time `json:"date"`
	Version  int       `json:"version"`
}

type WorkoutModel struct {
	DB *sql.DB
}

func (w WorkoutModel) Insert(workout *Workout) error {
	query := `
		INSERT INTO workouts (member_id, date)
		VALUES ($1, $2)
		RETURNING id, version
	`
	args := []interface{}{workout.MemberID, workout.Date}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return w.DB.QueryRowContext(ctx, query, args...).Scan(&workout.ID, &workout.Version)
}
