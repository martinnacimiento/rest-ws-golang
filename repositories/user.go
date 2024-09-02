package repositories

import (
	"context"

	"tincho.dev/rest-ws/models"
)

type UserRepository interface {
	FindAllUsers(ctx context.Context) ([]models.User, error)
	FindUserById(ctx context.Context, id int64) (*models.User, error)
	CreateUser(ctx context.Context, u *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

var implementation UserRepository

func SetUserRepository(repository UserRepository) {
	implementation = repository
}

func CreateUser(ctx context.Context, u *models.User) error {
	return implementation.CreateUser(ctx, u)
}

func FindAllUsers(ctx context.Context) ([]models.User, error) {
	return implementation.FindAllUsers(ctx)
}

func FindUserById(ctx context.Context, id int64) (*models.User, error) {
	return implementation.FindUserById(ctx, id)
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.GetUserByEmail(ctx, email)
}
