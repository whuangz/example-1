package repository

import (
	"context"
	"strconv"
	"time"

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

	query := "INSERT INTO post(title, content, author_id, created_at) VALUES (?, ?, ?, ?)"
	now := time.Now()
	result, err := p.db.ExecContext(ctx, query, post.Title, post.Content, post.AuthorID, now)
	if err != nil {
		return err
	}
	post.CreatedAt = now
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	post.ID = int32(id)
	post.AuthorID = 1
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

func (p *postRepo) FindByID(ctx context.Context, id int32) (domain.Post, error) {
	post := domain.Post{}
	query := "SELECT * FROM post WHERE id = ? LIMIT 1"
	rows, err := p.db.QueryContext(ctx, query, id)

	if err != nil {
		return post, err
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author.ID, &post.UpdatedAt, &post.CreatedAt)
		return post, nil
	} else {
		return post, domain.NewNotFound("id", strconv.Itoa(int(id)))
	}
}

func (p *postRepo) Update(ctx context.Context, id int32, post *domain.Post) error {
	query := `UPDATE post set title=?, content=?, updated_at=? WHERE id = ?`
	now := time.Now()
	//Coalesce
	res, err := p.db.ExecContext(ctx, query, post.Title, post.Content, now, post.ID)
	if err != nil {
		return err
	}

	post.UpdatedAt.Time = now

	if affect, _ := res.RowsAffected(); affect != 1 {
		return domain.NewBadRequest("id not found")
	}

	return nil
}

func (p *postRepo) Delete(ctx context.Context, id int32) error {
	query := "DELETE FROM post WHERE id = ?"
	results, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	if rowsAfected, _ := results.RowsAffected(); rowsAfected != 1 {
		return domain.NewBadRequest("id not found")
	}

	return nil
}
