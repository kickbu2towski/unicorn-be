package data

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"github.com/kickbu2towski/unicorn-be/cmd/api/internal/validator"
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

func ValidateToken(v *validator.Validator, token string) {
	v.Check(token == "", "token", "cannot be empty")
	v.Check(len(token) != 26, "token", "must be 26 characters long")
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

func (m *TokenModel) DeleteAllForUser(userID int64, scope string) error {
	query := `
	  DELETE FROM tokens
		WHERE user_id = $1 AND scope = $2;
	`
	_, err := m.DB.Exec(query, userID, scope)
	return err
}
