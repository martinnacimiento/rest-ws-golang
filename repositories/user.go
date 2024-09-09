package repositories

import (
	"context"

	"tincho.dev/rest-ws/models"
)

type UserRepository interface {
	FindAllUsers(ctx context.Context) ([]models.User, error)
	ListUsers(ctx context.Context, offset int64, limit int64) ([]models.User, error)
	FindUserById(ctx context.Context, id int64) (*models.User, error)
	CreateUser(ctx context.Context, u *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateOneUser(ctx context.Context, u *models.User) error
	DeleteOneUser(ctx context.Context, id int64) error
}

var userImplementation UserRepository

func SetUserRepository(repository UserRepository) {
	userImplementation = repository
}

func CreateUser(ctx context.Context, u *models.User) error {
	return userImplementation.CreateUser(ctx, u)
}

func FindAllUsers(ctx context.Context) ([]models.User, error) {
	return userImplementation.FindAllUsers(ctx)
}

func ListUsers(ctx context.Context, offset int64, limit int64) ([]models.User, error) {
	return userImplementation.ListUsers(ctx, offset, limit)
}

func FindUserById(ctx context.Context, id int64) (*models.User, error) {
	return userImplementation.FindUserById(ctx, id)
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return userImplementation.GetUserByEmail(ctx, email)
}

func UpdateOneUser(ctx context.Context, u *models.User) error {
	return userImplementation.UpdateOneUser(ctx, u)
}

func DeleteOneUser(ctx context.Context, id int64) error {
	return userImplementation.DeleteOneUser(ctx, id)
}
