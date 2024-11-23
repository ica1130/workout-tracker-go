package data

import (
	"context"
	"database/sql"
	"time"
)

type Member struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Height    int64     `json:"height"`
	Weight    int64     `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
}

type MemberModel struct {
	DB *sql.DB
}

func (m MemberModel) Insert(member *Member) error {
	query := `
		INSERT INTO members (email, name, height, weight)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	args := []interface{}{member.Email, member.Name, member.Height, member.Weight}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&member.ID, &member.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (m MemberModel) GetByEmail(email string) (*Member, error) {
	query := `
		SELECT id, email, name, height, weight, created_at
		FROM members
		WHERE email = $1
	`

	var member Member

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&member.ID,
		&member.Email,
		&member.Name,
		&member.Height,
		&member.Weight,
		&member.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &member, nil
}

func (m MemberModel) Update(member *Member) error {
	return nil
}

func (m MemberModel) Delete(id int64) error {
	return nil
}
