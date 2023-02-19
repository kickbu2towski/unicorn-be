package data

import (
	"database/sql"
	"errors"
)

var ErrNoRecord = errors.New("record not found")

type Models struct {
	TodoItems TodoItemModel
}

func NewModels(DB *sql.DB) Models {
	return Models{
		TodoItems: TodoItemModel{
			DB: DB,
		},
	}
}
