package data

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	Hash   []byte    `json:"hash"`
	Expiry time.Time `json:"expiry"`
	UserID int64     `json:"user_id"`
	Scope  string    `json:"scope"`
}

type TokenModel struct {
	DB *sql.DB
}

func generateToken() (string, []byte, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", nil, err
	}
	token := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes)
	hash := sha256.Sum256([]byte(token))
	return token, hash[:], nil
}

func (m *TokenModel) New(userID int64, ttl time.Duration, scope string) (string, error) {
	token, hash, err := generateToken()
	if err != nil {
		return "", err
	}

	t := &Token{
		Hash:   hash,
		Scope:  scope,
		Expiry: time.Now().Add(ttl),
		UserID: userID,
	}

	err = m.Insert(t)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (m *TokenModel) Insert(token *Token) error {
	query := `
	  INSERT INTO tokens(hash, expiry, scope, user_id)
		VALUES ($1, $2, $3, $4);
	`
	args := []any{token.Hash, token.Expiry, token.Scope, token.UserID}
	_, err := m.DB.Exec(query, args...)
	return err
}
