package store

import (
	"context"
	"fmt"

	"github.com/adocoder12/social_blog/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentStore struct {
	pool *pgxpool.Pool
}

func NewCommentStore(pool *pgxpool.Pool) *CommentStore {
	return &CommentStore{pool: pool}
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int64) ([]model.Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC`

	rows, err := s.pool.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("comment repo - get by post id query: %w", err)
	}
	defer rows.Close()

	comments := []model.Comment{}
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.Username); err != nil {
			return nil, fmt.Errorf("comment repo - get by post id scan: %w", err)
		}
		c.User.ID = c.UserID
		comments = append(comments, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("comment repo - get by post id rows: %w", err)
	}

	return comments, nil
}

func (s *CommentStore) Create(ctx context.Context, comment *model.Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err := s.pool.QueryRow(ctx, query, comment.PostID, comment.UserID, comment.Content).
		Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		if constraint, ok := isForeignKeyViolation(err); ok {
			switch constraint {
			case "comments_post_id_fkey":
				return ErrPostNotFound
			case "comments_user_id_fkey":
				return ErrUserNotFound
			}
		}
		return fmt.Errorf("comment repo - create: %w", err)
	}

	return nil
}
