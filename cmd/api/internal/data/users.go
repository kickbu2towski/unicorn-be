package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/kickbu2towski/unicorn-be/cmd/api/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicateEmail = errors.New("duplicate email")

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Activated bool      `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) setHash(password string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	u.Password = string(hash)
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name == "", "name", "cannot be empty")
	v.Check(len(user.Name) < 2, "name", "should be atleast 2 characters long")

	v.Check(user.Password == "", "password", "cannot be empty")
	v.Check(len(user.Password) < 6, "password", "should be atleast 6 characters long")
	v.Check(len(user.Password) > 72, "password", "should be maximum 72 characters long")

	v.Check(user.Email == "", "email", "cannot be empty")
	v.Check(!validator.Matches(user.Email, validator.EmailRegexp), "email", "invalid email address")
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(user *User) error {
	query := `
	  INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at;
	`
	args := []any{user.Name, user.Email, user.Password}
	err := m.DB.QueryRow(query, args...).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}
