package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/whuangz/go-example/go-api/domain"
)

type postRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) domain.PostRepository {
	return &postRepo{db: db}
}

func (p *postRepo) Save(ctx context.Context, post *domain.Post) error {

	query := "INSERT INTO post(title, content, author_id) VALUES (?, ?, ?)"
	result, err := p.db.ExecContext(ctx, query, post.Title, post.Content, 1)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	post.ID = int32(id)
	post.Author.ID = 1
	return nil
}

func (p *postRepo) FindAll(ctx context.Context) ([]domain.Post, error) {

	posts := []domain.Post{}
	query := "SELECT * FROM post"
	rows, err := p.db.QueryContext(ctx, query)

	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		post := domain.Post{}
		rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author.ID, &post.UpdatedAt, &post.CreatedAt)
		posts = append(posts, post)
	}
	return posts, nil
}

func (p *postRepo) FindByID(ctx context.Context, id uint32) (domain.Post, error) {
	return domain.Post{}, nil
}

func (p *postRepo) Update(ctx context.Context, id uint32, post *domain.Post) error {
	return nil
}

func (p *postRepo) Delete(ctx context.Context, id uint32) (int32, error) {
	return 0, nil
}
