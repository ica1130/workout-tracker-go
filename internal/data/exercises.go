package data

import (
	"context"
	"database/sql"
	"time"
)

type Exercise struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Version     int    `json:"version"`
}

type ExerciseModel struct {
	DB *sql.DB
}

func (e ExerciseModel) Insert(exercise *Exercise) error {
	query := `
		INSERT INTO exercises (name, category, description)
		VALUES ($1, $2, $3)
		RETURNING id
		`

	args := []interface{}{exercise.Name, exercise.Category, exercise.Description}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, args...).Scan(&exercise.ID)
	if err != nil {
		return err
	}

	return nil
}
