package domains

import (
	"database/sql"

	"github.com/go-playground/validator/v10"
)

type Post struct {
	ID        int32        `json:id`
	Title     string       `json:title validate:"required"`
	Content   string       `json:content validate:"required"`
	Password  string       `json:password validate:"required"`
	UpdatedAt sql.NullTime `json:updated_at`
	CreatedAt string       `json:created_at`
	Author    Author       `json:author`
}

func (p *Post) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
