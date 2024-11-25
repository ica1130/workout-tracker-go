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

	return e.DB.QueryRowContext(ctx, query, args...).Scan(&exercise.ID, &exercise.Version)
}

func (e ExerciseModel) GetByCategory(category string) ([]*Exercise, error) {
	query := `
		SELECT id, name, category, description, version
		FROM exercises
		WHERE category = $1
		ORDER BY id`

	args := []interface{}{category}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := e.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	exercises := []*Exercise{}

	for rows.Next() {
		var exercise Exercise

		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Category,
			&exercise.Description,
			&exercise.Version,
		)

		if err != nil {
			return nil, err
		}

		exercises = append(exercises, &exercise)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return exercises, nil
}
