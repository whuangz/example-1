package service

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/whuangz/go-example/go-api/domain"
	"github.com/whuangz/go-example/go-api/helpers/jwt"
)

type tokenService struct {
	repo            domain.TokenRepository
	privKey         *rsa.PrivateKey
	pubKey          *rsa.PublicKey
	refreshSecret   string
	accessTokenExp  int64
	refreshTokenExp int64
}

func NewTokenService(repo domain.TokenRepository, private *rsa.PrivateKey, public *rsa.PublicKey, refresh string, accTokenExp int64, refreshTokenExp int64) domain.TokenService {
	return &tokenService{repo, private, public, refresh, accTokenExp, refreshTokenExp}
}

func (s *tokenService) NewPairFromUser(ctx context.Context, a *domain.Account, prevAccesstoken string) (*domain.TokenPair, error) {

	if prevAccesstoken != "" {
		if err := s.repo.DeleteRefreshToken(ctx, a.UID.String(), prevAccesstoken); err != nil {
			log.Printf("Could not delete previous refreshToken for uid: %v, tokenID: %v\n", a.UID.String(), prevAccesstoken)

			return nil, err
		}
	}

	accToken, err := jwt.GenerateAccessToken(a, s.privKey, s.accessTokenExp)
	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", a.UID, err.Error())
		return nil, domain.NewInternal()
	}

	refreshToken, err := jwt.GenerateRefreshToken(a.UID, s.refreshSecret, s.refreshTokenExp)

	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", a.UID, err.Error())
		return nil, domain.NewInternal()
	}

	if err := s.repo.SetRefreshToken(ctx, a.UID.String(), refreshToken.ID, refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", a.UID, err.Error())
		return nil, domain.NewInternal()
	}

	return &domain.TokenPair{AccessToken: accToken, RefreshToken: refreshToken.SS}, nil

}

func (s *tokenService) ValidateAccessToken(token string) (*domain.Account, error) {
	claims, err := jwt.ValidateAccessToken(token, s.pubKey)

	if err != nil {
		return nil, domain.NewAuthorization("Unable to verify user")
	}

	return claims.Account, nil
}
