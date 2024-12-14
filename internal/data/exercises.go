package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"workout-tracker-go.ilijakrilovic.com/internal/validator"
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

type Exercise struct {
	ID          int64  `json:"-"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Version     int    `json:"-"`
}

func ValidateCategory(v *validator.Validator, category string) {
	v.Check(category != "", "category", "must be provided")
	v.Check(allowedCategories[category], "category", "invalid category")
}

func ValidateExercise(v *validator.Validator, exercise *Exercise) {
	v.Check(exercise.Name != "", "name", "must be provided")
	v.Check(len(exercise.Name) <= 50, "name", "must not be more than 50 butes long")

	ValidateCategory(v, exercise.Category)

	v.Check(len(exercise.Description) <= 1000, "description", "must not be more than 1000 bytes long")
}

type ExerciseModel struct {
	DB *sql.DB
}

func (e ExerciseModel) Insert(exercise *Exercise) error {
	query := `
		INSERT INTO exercises (name, category, description)
		VALUES ($1, $2, $3)
		RETURNING id, version
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

func (e ExerciseModel) GetById(id int64) (*Exercise, error) {
	query := `
		SELECT id, name, category, description, version
		FROM exercises
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exercise Exercise

	err := e.DB.QueryRowContext(ctx, query, id).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Category,
		&exercise.Description,
		&exercise.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &exercise, nil
}

func (e ExerciseModel) Update(exercise *Exercise) error {
	query := `
		UPDATE exercises
		SET name = $1, category = $2, description = $3, version = version +1
		WHERE id = $4 AND version = $5
		RETURNING version
		`

	args := []interface{}{exercise.Name, exercise.Category, exercise.Description, exercise.ID, exercise.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, args...).Scan(&exercise.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (e ExerciseModel) Delete(id int64) error {
	query := `
		DELETE FROM exercises
		WHERE id = $1
	`
	result, err := e.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
