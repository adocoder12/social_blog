package store

import (
	"context"

	"github.com/adocoder12/social_blog/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdatePassword(ctx context.Context, id int64, hash []byte) error
	DeleteUser(ctx context.Context, id int64) error
}

type PostRepository interface {
	GetAllPosts(ctx context.Context) ([]model.Post, error)
	GetPostByID(ctx context.Context, id int64) (*model.Post, error)
	CreatePost(ctx context.Context, post *model.Post) (*model.Post, error)
	UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error)
	DeletePost(ctx context.Context, id int64) error
}

type CommentRepository interface {
	GetByPostID(ctx context.Context, id int64) ([]model.Comment, error)
	Create(ctx context.Context, comment *model.Comment) error 
}

type Storage struct {
	Users UserRepository
	Posts PostRepository
	Comments CommentRepository
}

func NewPostgresStorage(pool *pgxpool.Pool) Storage {
	return Storage{
		Users: &UsersStore{pool: pool},
		Posts: &PostStore{pool: pool},
		Comments: &CommentStore{pool: pool},
	}
}
