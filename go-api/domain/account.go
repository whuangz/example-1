package domain

import (
	"context"

	"github.com/google/uuid"
)

type AccountRepository interface {
	FindByID(ctx context.Context, uid uuid.UUID) (*Account, error)
}

type AccountService interface {
	Get(ctx context.Context, uid uuid.UUID) (*Account, error)
	Signup(ctx context.Context, a *Account) error
}

type Account struct {
	UID      uuid.UUID `json:"uid"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Name     string    `json:"name"`
	ImageUrl string    `json:"image_url"`
	Website  string    `json:"website"`
}
