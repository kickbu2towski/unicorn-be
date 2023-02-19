package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/kickbu2towski/unicorn-be/cmd/api/internal/validator"
	"github.com/lib/pq"
)

type TodoItemModel struct {
	DB *sql.DB
}

type TodoItem struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	State     string     `json:"state"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
	Tags      []string   `json:"tags"`
	Priority  string     `json:"priority"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func ValidateTodoItem(v *validator.Validator, todoItem *TodoItem) {
	v.Check(todoItem.Name == "", "name", "cannot be empty")
	v.Check(todoItem.State == "", "state", "cannot be empty")
}

func (m *TodoItemModel) Insert(todoItem *TodoItem) error {
	query := `
	  INSERT INTO todo_items(name, state, priority, tags)
		VALUES($1, $2, $3, $4)
		RETURNING id, created_at;
	`
	args := []any{todoItem.Name, todoItem.State, todoItem.Priority, pq.Array(todoItem.Tags)}
	return m.DB.QueryRow(query, args...).Scan(&todoItem.Id, &todoItem.CreatedAt)
}

func (m *TodoItemModel) Get(id int64) (*TodoItem, error) {
	query := `
	  SELECT id, state, priority, tags, closed_at, created_at
		FROM todo_items
		WHERE id = $1;
	`

	var todoItem TodoItem

	err := m.DB.QueryRow(query, id).Scan(
		&todoItem.Id,
		&todoItem.State,
		&todoItem.Priority,
		pq.Array(&todoItem.Tags),
		&todoItem.ClosedAt,
		&todoItem.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecord
		default:
			return nil, err
		}
	}

	return &todoItem, nil
}

func (m *TodoItemModel) Delete(id int64) error {
	query := `
	  DELETE FROM todo_items
		WHERE id = $1;
	`

	res, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNoRecord
	}

	return nil
}

func (m *TodoItemModel) Update(id int64, todoItem *TodoItem) error {
	query := `
	  UPDATE todo_items
		SET name = $1, state = $2, priority = $3, tags = $4, closed_at = $5
		WHERE id = $6
	`
	args := []any{todoItem.Name, todoItem.State, todoItem.Priority, pq.Array(todoItem.Tags), todoItem.ClosedAt, id}

	_, err := m.DB.Exec(query, args...)
	return err
}

/*
    ** TODO [#A] Buy Milk :home:work:
		Closed: [2023-02-18]
*/
