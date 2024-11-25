package data

import "database/sql"

type Exercise struct {
	ID          int64  `json:"id"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type ExerciseModel struct {
	DB *sql.DB
}
