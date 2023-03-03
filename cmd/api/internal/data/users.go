package data

import (
	"crypto/sha256"
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

func (u *User) GenerateHash() {
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	u.Password = string(hash)
}

func (u *User) PasswordMatches(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email == "", "email", "cannot be emty")
	v.Check(!validator.Matches(email, validator.EmailRegexp), "email", "invalid email address")
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password == "", "password", "cannot be empty")
	v.Check(len(password) < 6, "password", "should be atleast 6 characters long")
	v.Check(len(password) > 72, "password", "should be maximum 72 characters long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name == "", "name", "cannot be empty")
	v.Check(len(user.Name) < 2, "name", "should be atleast 2 characters long")
	ValidatePassword(v, user.Password)
	ValidateEmail(v, user.Email)
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

func (m *UserModel) GetForToken(token string, scope string) (*User, error) {
	var user User
	hash := sha256.Sum256([]byte(token))

	query := `
	  SELECT u.id, u.name, u.email, u.password, u.activated, u.created_at
		FROM users u
		INNER JOIN tokens t
		ON u.id = t.user_id
		WHERE u.activated = false AND t.hash = $1 AND t.expiry > $2 AND scope = $3;
	`

	args := []any{hash[:], time.Now(), scope}

	err := m.DB.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Activated,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecord
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m *UserModel) Update(user *User) error {
	query := `
	  UPDATE users 
		SET name = $1, activated = $2
		WHERE id = $3;
	`
	args := []any{user.Name, user.Activated, user.ID}
	_, err := m.DB.Exec(query, args...)
	return err
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	var user User

	query := `
	  SELECT id, password
		FROM users
		WHERE email = $1;
	`

	err := m.DB.QueryRow(query, email).Scan(&user.ID, &user.Password)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecord
		default:
			return nil, err
		}
	}

	return &user, nil
}
