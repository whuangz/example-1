package service

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/whuangz/go-example/go-api/domain"
	"github.com/whuangz/go-example/go-api/helpers/crypto"
	"github.com/whuangz/go-example/go-api/helpers/jwt"
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

	hash, err := crypto.HashPassword(a.Password)
	if err != nil {
		return domain.NewInternal()
	}

	a.Password = hash

	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.repo.Create(c, a); err != nil {
		return err
	}
	// If we get around to adding events, we'd Publish it here
	// err := s.EventsBroker.PublishAccountUpdated(a, true)

	// if err != nil {
	//  return nil, domain.NewInternal()
	// }

	return nil
}

func (s *accountService) Signin(ctx context.Context, a *domain.Account) error {
	accFetched, err := s.repo.FindByEmail(ctx, a.Email)

	// Will return NotAuthorized to client to omit details of why
	if err != nil {
		return domain.NewAuthorization("Invalid email and password combination")
	}

	// verify password - we previously created this method
	match, err := crypto.ValidateHash(accFetched.Password, a.Password)

	if err != nil {
		return domain.NewInternal()
	}

	if !match {
		return domain.NewAuthorization("Invalid email and password combination")
	}

	*a = *accFetched
	return nil
}

func (s *tokenService) ValidateRefreshToken(tokenString string) (*domain.RefreshToken, error) {
	// validate actual JWT with string a secret
	claims, err := jwt.ValidateRefreshToken(tokenString, s.refreshSecret)

	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		log.Printf("Unable to validate or parse refreshToken for token string: %s\n%v\n", tokenString, err)
		return nil, domain.NewAuthorization("Unable to verify user from refresh token")
	}

	// Standard claims store ID as a string. I want "model" to be clear our string
	// is a UUID. So we parse claims.Id as UUID
	tokenUUID, err := uuid.Parse(claims.Id)

	if err != nil {
		log.Printf("Claims ID could not be parsed as UUID: %s\n%v\n", claims.Id, err)
		return nil, domain.NewAuthorization("Unable to verify user from refresh token")
	}

	return &domain.RefreshToken{
		SS:  tokenString,
		ID:  tokenUUID,
		UID: claims.UID,
	}, nil
}
