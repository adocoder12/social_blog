package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/adocoder12/social_blog/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type PostStore struct {
	pool *pgxpool.Pool
}

func (s *PostStore) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	query := `SELECT id, title, content, user_id, tags, created_at, updated_at FROM posts`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("post repo - get all query: %w", err)
	}
	defer rows.Close()

	posts := []model.Post{}
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Tags, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, fmt.Errorf("post repo - get all scan: %w", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("post repo - get all rows: %w", err)
	}
	return posts, nil
}

func (s *PostStore) GetPostByID(ctx context.Context, id int64) (*model.Post, error) {
	query := `SELECT id, title, content, user_id, tags, created_at, updated_at FROM posts WHERE id = $1`
	var post model.Post
	err := s.pool.QueryRow(ctx, query, id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Tags, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, fmt.Errorf("post repo - get by id query row: %w", err)
	}
	return &post, nil
}

func (s *PostStore) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	query := `
		INSERT INTO posts (title, content, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`
	// WE DONT NEED TO USE pq.Array CUZ pgx v5 handles Go slices natively
	if post.Tags == nil {
		post.Tags = []string{} // nil would be sent as NULL and violate NOT NULL
	}

	err := s.pool.QueryRow(ctx, query, post.Title, post.Content, post.UserID, post.Tags).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("post repo - create post query: %w", err)
	}
	return post, nil
}

func (s *PostStore) UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	query := `
		UPDATE posts
		SET title = $1, content = $2, tags = $3, updated_at = now()
		WHERE id = $4
		RETURNING user_id, created_at, updated_at`

	if post.Tags == nil {
		post.Tags = []string{}
	}

	err := s.pool.QueryRow(ctx, query, post.Title, post.Content, post.Tags, post.ID).
		Scan(&post.UserID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, fmt.Errorf("post repo - update query: %w", err)
	}
	return post, nil
}

func (s *PostStore) DeletePost(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`
	result, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("post repo - delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrPostNotFound
	}
	return nil
}
