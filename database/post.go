package database

import (
	"context"

	"tincho.dev/rest-ws/models"
)

func (p *Postgres) CreatePost(ctx context.Context, post *models.Post) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO posts (title, content, user_id) VALUES ($1, $2, $3)", post.Title, post.Content, post.UserID)

	return err
}

func (p *Postgres) FindAllPosts(ctx context.Context) ([]models.Post, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, title, content, user_id, created_at, updated_at FROM posts")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []models.Post{}

	for rows.Next() {
		var post models.Post

		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (p *Postgres) FindPostById(ctx context.Context, id int64) (*models.Post, error) {
	row := p.db.QueryRowContext(ctx, "SELECT id, title, content, user_id, created_at, updated_at FROM posts WHERE id = $1", id)

	var post models.Post

	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *Postgres) FindPostsByUserId(ctx context.Context, userId int64) ([]models.Post, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, title, content, user_id, created_at, updated_at FROM posts WHERE user_id = $1", userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []models.Post{}

	for rows.Next() {
		var post models.Post

		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (p *Postgres) UpdatePost(ctx context.Context, post *models.Post) error {
	_, err := p.db.ExecContext(ctx, "UPDATE posts SET title = $1, content = $2, updated_at = NOW() WHERE id = $3", post.Title, post.Content, post.ID)

	return err
}

func (p *Postgres) DeletePost(ctx context.Context, id int64) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1", id)

	return err
}

func (p *Postgres) Close() error {
	return p.db.Close()
}
