package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/adocoder12/social_blog/internal/model"
	"github.com/jackc/pgx/v5"
)

var ErrPostNotFound = errors.New("post not found")

type PostStore struct {
	db DBTX
}

func NewPostStore(db DBTX) *PostStore {
	return &PostStore{db: db}
}

func (s *PostStore) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	query := `
		SELECT id, title, content, user_id, tags, version ,created_at, updated_at
		FROM posts
		ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("post repo - get all query: %w", err)
	}
	defer rows.Close()

	posts := []model.Post{}
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Tags,&post.Version ,&post.CreatedAt, &post.UpdatedAt); err != nil {
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
	query := `
		SELECT id, title, content, user_id, tags, version, created_at, updated_at
		FROM posts
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post model.Post
	err := s.db.QueryRow(ctx, query, id).
		Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Tags, &post.Version,&post.CreatedAt, &post.UpdatedAt)
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
	// pgx v5 handles Go slices natively; nil would be sent as NULL and violate NOT NULL.
	if post.Tags == nil {
		post.Tags = []string{}
	}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRow(ctx, query, post.Title, post.Content, post.UserID, post.Tags).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if _, ok := isForeignKeyViolation(err); ok {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("post repo - create post query: %w", err)
	}
	return post, nil
}

func (s *PostStore) UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	query := `
		UPDATE posts
		SET title = $1, content = $2, tags = $3, updated_at = now() , version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING user_id, version, updated_at`

	if post.Tags == nil {
		post.Tags = []string{}
	}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRow(ctx, query, post.Title, post.Content, post.Tags, post.ID,post.Version).
		Scan(&post.UserID, &post.Version, &post.UpdatedAt)
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("post repo - delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrPostNotFound
	}
	return nil
}
