package data

import (
	"database/sql"
	"time"
)

type Member struct {
	ID        int64
	Email     string
	Name      string
	Height    int64
	Weight    int64
	CreatedAt time.Time
}

type MemberModel struct {
	DB *sql.DB
}

func (m MemberModel) Insert(member *Member) error {
	return nil
}

func (m MemberModel) Get(id int64) (*Member, error) {
	return nil, nil
}

func (m MemberModel) Update(member *Member) error {
	return nil
}

func (m MemberModel) Delete(id int64) error {
	return nil
}
