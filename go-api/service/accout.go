package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/whuangz/go-example/go-api/domain"
)

type accountService struct {
	repo domain.AccountRepository
}

func NewAccountService(repo domain.AccountRepository) domain.AccountService {
	return &accountService{repo: repo}
}

func (s *accountService) Get(ctx context.Context, uid uuid.UUID) (*domain.Account, error) {
	a, err := s.repo.FindByID(ctx, uid)
	return a, err
}

func (s *accountService) Signup(ctx context.Context, a *domain.Account) error {
	panic("Method not implemented")
}
