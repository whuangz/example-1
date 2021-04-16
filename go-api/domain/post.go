package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
)

type PostRepository interface {
	Save(ctx context.Context, post *Post) error
	FindAll(ctx context.Context) ([]Post, error)
	FindByID(ctx context.Context, id int32) (Post, error)
	Update(ctx context.Context, id int32, post *Post) error
	Delete(ctx context.Context, id int32) error
}

type PostService interface {
	Save(ctx context.Context, post *Post) error
	FindAll(ctx context.Context) ([]Post, error)
	FindByID(ctx context.Context, id int32) (Post, error)
	Update(ctx context.Context, id int32, post *Post) error
	Delete(ctx context.Context, id int32) error
}

type Post struct {
	ID        int32        `json:id valid:"omitempty"`
	Title     string       `json:title valid:"omitempty"`
	Content   string       `json:content valid:"omitempty"`
	UpdatedAt sql.NullTime `json:updated_at`
	CreatedAt time.Time    `json:created_at`
	Author    Author       `json:author`
	AuthorID  int32        `json:author_id`
}

func (p *Post) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
