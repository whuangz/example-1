package domain

import (
	"database/sql"

	"github.com/go-playground/validator/v10"
)

type Author struct {
	ID        int32        `json:id`
	Username  string       `json:username validate:"required"`
	Email     string       `json:email validate:"required"`
	UpdatedAt sql.NullTime `json:updated_at`
	CreatedAt string       `json:created_at`
}

func (a *Author) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}
