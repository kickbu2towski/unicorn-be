package data

import (
	"database/sql"
	"errors"
)

var ErrNoRecord = errors.New("record not found")

type Models struct {
	UserModel UserModel
}

func NewModels(DB *sql.DB) Models {
	return Models{
		UserModel: UserModel{
			DB: DB,
		},
	}
}
