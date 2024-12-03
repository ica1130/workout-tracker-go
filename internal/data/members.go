package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Member struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Height    int64     `json:"height"`
	Weight    int64     `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plain string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plain
	p.hash = hash

	return nil
}

func (p *password) Compare() (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(*p.plaintext))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type MemberModel struct {
	DB *sql.DB
}

func (m MemberModel) Insert(member *Member) error {
	query := `
		INSERT INTO members (email, name, password_hash, activated, height, weight)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, version
	`

	args := []interface{}{member.Email, member.Name, member.Password.hash, member.Activated, member.Height, member.Weight}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&member.ID, &member.CreatedAt, &member.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "members_email_key"`:
			return errors.New("duplicate email")
		default:
			return err
		}
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

func (m MemberModel) GetForToken(tokenScope, tokenPlain string) (*Member, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlain))

	query := `
		SELECT m.id, m.email, m.name, m.password_hash, m.activated, m.height, m.weight, m.created_at, m.version
		FROM members m
		INNER JOIN tokens
		ON m.id = tokens.member_id
		WHERE tokens.hash = $1
		AND tokens.scope = $2
		AND tokens.expiry > $3 
	`

	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var member Member

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&member.ID,
		&member.Email,
		&member.Name,
		&member.Password.hash,
		&member.Activated,
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

/* 	ID        int64     `json:"id"`
Email     string    `json:"email"`
Name      string    `json:"name"`
Password  password  `json:"-"`
Activated bool      `json:"activated"`
Height    int64     `json:"height"`
Weight    int64     `json:"weight"`
CreatedAt time.Time `json:"created_at"`
Version   int       `json:"-"` */
