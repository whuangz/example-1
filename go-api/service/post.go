package service

import (
	"context"
	"strconv"
	"time"

	"github.com/whuangz/go-example/go-api/domain"
)

type postService struct {
	repo domain.PostRepository
}

func NewPostService(repo domain.PostRepository) domain.PostService {
	return &postService{repo: repo}
}

func (p *postService) Save(ctx context.Context, post *domain.Post) error {
	err := post.Validate()
	if err != nil {
		return err
	}
	post.AuthorID = 1
	return p.repo.Save(ctx, post)
}

func (p *postService) FindAll(ctx context.Context) ([]domain.Post, error) {
	posts, err := p.repo.FindAll(ctx)
	return posts, err
}

func (p *postService) FindByID(ctx context.Context, id int32) (domain.Post, error) {
	post, err := p.repo.FindByID(ctx, id)
	return post, err
}

func (p *postService) Update(ctx context.Context, id int32, post *domain.Post) error {
	if id == 0 {
		return domain.NewNotFound("id", strconv.Itoa(int(id)))
	}
	post.UpdatedAt.Time = time.Now()
	return p.repo.Update(ctx, id, post)
}

func (p *postService) Delete(ctx context.Context, id int32) error {
	if id == 0 {
		return domain.NewNotFound("id", strconv.Itoa(int(id)))
	}
	return p.repo.Delete(ctx, id)
}
