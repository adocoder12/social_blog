package store

import (
	"context"

	"github.com/adocoder12/social_blog/internal/model"
)

type UserRepositoryInterface interface {
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdatePassword(ctx context.Context, id int, hash string) error
	DeleteUser(ctx context.Context, id int) error
}
