package store

import "github.com/jackc/pgx/v5/pgxpool"

type PostStore struct {
	pool *pgxpool.Pool
}

func (s *PostStore) CreatePost() {

}
