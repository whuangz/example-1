package service

import (
	"context"

	"github.com/whuangz/go-example/go-api/domain"
)

type postService struct {
	repo domain.PostRepository
}

func NewPostService(repo domain.PostRepository) domain.PostService {
	return &postService{repo: repo}
}

func (p *postService) FindAll(ctx context.Context) ([]domain.Post, error) {
	posts, err := p.repo.FindAll(ctx)
	return posts, err
}
