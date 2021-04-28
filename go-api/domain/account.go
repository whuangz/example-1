package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type AccountRepository interface {
	FindByEmail(ctx context.Context, email string) (*Account, error)
	FindByID(ctx context.Context, uid uuid.UUID) (*Account, error)
	Create(ctx context.Context, a *Account) error
}

type AccountService interface {
	Get(ctx context.Context, uid uuid.UUID) (*Account, error)
	Signup(ctx context.Context, a *Account) error
	Signin(ctx context.Context, a *Account) error
}

type Account struct {
	ID        int32        `json:"id"`
	UID       uuid.UUID    `json:"uid"`
	Email     string       `json:"email"`
	Password  string       `json:"-"`
	Name      string       `json:"name"`
	ImageUrl  string       `json:"image_url"`
	Website   string       `json:"website"`
	UpdatedAt sql.NullTime `json:updated_at`
	CreatedAt time.Time    `json:created_at`
}

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, accID string, accessToken string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, accID string, prevAccessToken string) error
}

type TokenService interface {
	NewPairFromUser(ctx context.Context, a *Account, prevAccesstoken string) (*TokenPair, error)
	ValidateAccessToken(token string) (*Account, error)
	ValidateRefreshToken(tokenString string) (*RefreshToken, error)
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	ID  uuid.UUID `json:"-"`
	UID uuid.UUID `json:"-"`
	SS  string    `json:"refresh_token"`
}
