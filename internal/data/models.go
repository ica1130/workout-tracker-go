package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Members MemberModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Members: MemberModel{DB: db},
	}
}
