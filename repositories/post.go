package repositories

import (
	"context"

	"tincho.dev/rest-ws/models"
)

type PostRepository interface {
	FindAllPosts(ctx context.Context) ([]models.Post, error)
	FindPostById(ctx context.Context, id int64) (*models.Post, error)
	CreatePost(ctx context.Context, p *models.Post) error
	UpdatePost(ctx context.Context, p *models.Post) error
	DeletePost(ctx context.Context, id int64) error
}

var postImplementation PostRepository

func SetPostRepository(repository PostRepository) {
	postImplementation = repository
}

func FindAllPosts(ctx context.Context) ([]models.Post, error) {
	return postImplementation.FindAllPosts(ctx)
}

func FindPostById(ctx context.Context, id int64) (*models.Post, error) {
	return postImplementation.FindPostById(ctx, id)
}

func CreatePost(ctx context.Context, p *models.Post) error {
	return postImplementation.CreatePost(ctx, p)
}

func UpdatePost(ctx context.Context, p *models.Post) error {
	return postImplementation.UpdatePost(ctx, p)
}

func DeletePost(ctx context.Context, id int64) error {
	return postImplementation.DeletePost(ctx, id)
}
