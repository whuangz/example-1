package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/whuangz/go-example/go-api/domain"
)

type accountRepo struct {
	db *sqlx.DB
}

func NewAccountRepo(db *sqlx.DB) domain.AccountRepository {
	return &accountRepo{db: db}
}

func (r *accountRepo) FindByEmail(ctx context.Context, email string) (*domain.Account, error) {
	acc := &domain.Account{}
	query := "SELECT id, uid, email, name, image_url, website FROM account WHERE email=?"
	rows, err := r.db.QueryContext(ctx, query, email)

	if err != nil {
		return acc, err
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&acc.ID, &acc.UID, &acc.Email, &acc.Name, &acc.ImageUrl, &acc.Website)
		return acc, nil
	} else {
		return acc, domain.NewNotFound("email", email)
	}
}

func (r *accountRepo) FindByID(ctx context.Context, uid uuid.UUID) (*domain.Account, error) {

	acc := &domain.Account{}
	query := "SELECT id, uid, email, name, password, image_url, website FROM account WHERE uid=?"
	rows, err := r.db.QueryContext(ctx, query, uid)

	if err != nil {
		return acc, err
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&acc.ID, &acc.UID, &acc.Email, &acc.Name, &acc.Password, &acc.ImageUrl, &acc.Website)
		return acc, nil
	} else {
		return acc, domain.NewNotFound("uid", uid.String())
	}
}

func (r *accountRepo) Create(ctx context.Context, a *domain.Account) error {
	query := "INSERT INTO account (email, password, uid, created_at) VALUES (?, ?, ?, ?)"
	uid := uuid.New()
	now := time.Now()

	result, err := r.db.ExecContext(ctx, query, a.Email, a.Password, uid, now)
	if err != nil {
		return err
	}
	a.CreatedAt = now
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	a.ID = int32(id)
	a.UID = uid

	return nil

}
