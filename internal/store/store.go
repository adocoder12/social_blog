package store

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Users UserRepositoryInterface
	Posts interface{ any }
}

func NewPostgresStorage(pool *pgxpool.Pool) Storage {
	return Storage{
		Users: &UsersStore{pool},
		Posts: &PostStore{pool},
	}
}
