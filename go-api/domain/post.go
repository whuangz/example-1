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
	FindByID(ctx context.Context, id uint32) (Post, error)
	Update(ctx context.Context, id uint32, post *Post) error
	Delete(ctx context.Context, id uint32) (int32, error)
}

type PostService interface {
	FindAll(ctx context.Context) ([]Post, error)
}

type Post struct {
	ID        int32        `json:id`
	Title     string       `json:title validate:"required"`
	Content   string       `json:content validate:"required"`
	UpdatedAt sql.NullTime `json:updated_at`
	CreatedAt time.Time    `json:created_at`
	Author    Author       `json:author`
}

func (p *Post) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
