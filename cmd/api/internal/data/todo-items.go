package data

import (
	"time"

	"github.com/kickbu2towski/unicorn-be/cmd/api/internal/validator"
)

type TodoItem struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	State    string     `json:"state"`
	ClosedAt *time.Time `json:"closed_at,omitempty"`
	Tags     []string   `json:"tags"`
	Priority string     `json:"priority"`
}

func ValidateTodoItem(v *validator.Validator, todoItem *TodoItem) {
	v.Check(todoItem.Name == "", "name", "cannot be empty")
	v.Check(todoItem.State == "", "state", "cannot be empty")
}

/*
    ** TODO [#A] Buy Milk :home:work:
		Closed: [2023-02-18]
*/
