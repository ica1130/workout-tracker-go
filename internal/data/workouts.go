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

type WorkoutResponse struct {
	ID       int64                    `json:"id"`
	MemberID int64                    `json:"-"`
	Date     time.Time                `json:"date"`
	Details  []*WorkoutDetailResponse `json:"details"`
	Version  int                      `json:"-"`
}

type WorkoutDetailResponse struct {
	ID          int64    `json:"-"`
	WorkoutID   int64    `json:"-"`
	Exercise    Exercise `json:"exercise"`
	Set         int      `json:"set"`
	Repetitions int      `json:"repetitions"`
	Weight      float64  `json:"weight"`
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

func (w WorkoutModel) GetByMemberID(memberID int64) ([]*WorkoutResponse, error) {
	query := `
		SELECT 	w.id, w.date, wd.set, wd.repetitions, wd.weight, e.name, e.category, e.description 
		FROM workouts w 
		JOIN workout_details AS wd 
		ON w.id = wd.workout_id
		JOIN exercises e
		ON e.id = wd.exercise_id
		WHERE w.member_id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := w.DB.QueryContext(ctx, query, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workouts := make(map[int64]*WorkoutResponse)

	for rows.Next() {
		var workout WorkoutResponse
		var detail WorkoutDetailResponse

		err := rows.Scan(
			&workout.ID,
			&workout.Date,
			&detail.Set,
			&detail.Repetitions,
			&detail.Weight,
			&detail.Exercise.Name,
			&detail.Exercise.Category,
			&detail.Exercise.Description,
		)

		if err != nil {
			return nil, err
		}

		if _, ok := workouts[workout.ID]; !ok {
			workout.Details = append(workout.Details, &detail)
			workouts[workout.ID] = &workout
		} else {
			updateWorkout := workouts[workout.ID]
			updateWorkout.Details = append(updateWorkout.Details, &detail)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	workoutsSlice := make([]*WorkoutResponse, 0, len(workouts))
	for _, workout := range workouts {
		fmt.Println(workout.ID)
		workoutsSlice = append(workoutsSlice, workout)
	}

	return workoutsSlice, nil
}

func (w WorkoutModel) Delete(id int64) error {
	query := `
		DELETE FROM workouts
		WHERE id = $1
	`
	result, err := w.DB.Exec(query, id)
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
