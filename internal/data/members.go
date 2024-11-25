package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Member struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Height    int64     `json:"height"`
	Weight    int64     `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"-"`
}

type MemberModel struct {
	DB *sql.DB
}

func (m MemberModel) Insert(member *Member) error {
	query := `
		INSERT INTO members (email, name, height, weight)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version
	`

	args := []interface{}{member.Email, member.Name, member.Height, member.Weight}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&member.ID, &member.CreatedAt, &member.Version)
	if err != nil {
		return err
	}

	return nil
}

func (m MemberModel) GetByEmail(email string) (*Member, error) {
	query := `
		SELECT id, email, name, height, weight, created_at, version
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
		&member.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &member, nil
}

func (m MemberModel) GetById(id int64) (*Member, error) {
	query := `
	SELECT id, email, name, height, weight, created_at, version
	FROM members
	WHERE id = $1
	`

	var member Member

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&member.ID,
		&member.Email,
		&member.Name,
		&member.Height,
		&member.Weight,
		&member.CreatedAt,
		&member.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &member, nil
}

func (m MemberModel) Update(member *Member) error {
	query := `
		UPDATE members
		SET email = $1, name = $2, height = $3, weight = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version
	`

	args := []interface{}{
		member.Email,
		member.Name,
		member.Height,
		member.Weight,
		member.ID,
		member.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&member.Version)
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

func (m MemberModel) Delete(id int64) error {

	query := `
		DELETE FROM members
		WHERE id = $1
	`

	result, err := m.DB.Exec(query, id)
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
