package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"
)

type Token struct {
	Plaintext string
	Hash      []byte
	MemberID  int64
	Expiry    time.Time
	Scope     string
}

func generateToken(memberID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		MemberID: memberID,
		Expiry:   time.Now().Add(ttl),
		Scope:    scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

type TokenModel struct {
	DB *sql.DB
}

func (m TokenModel) New(memberID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(memberID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

func (m TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO tokens (hash, member_id, expiry, scope)
		VALUES ($1, $2, $3, $4)
	`

	args := []interface{}{token.Hash, token.MemberID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}
