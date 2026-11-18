package store

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrAlreadyFollowing = errors.New("already following")
	ErrNotFollowing     = errors.New("not following this user")
)

type FollowerStore struct {
	db DBTX
}

func (s *FollowerStore) Follow(ctx context.Context, followerID, userID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.Exec(ctx, query, userID, followerID)
	if err != nil {
		if _, ok := isUniqueViolation(err); ok {
			return ErrAlreadyFollowing
		}
		if _, ok := isForeignKeyViolation(err); ok {
			return ErrUserNotFound
		}
		return fmt.Errorf("follower repo - follow: %w", err)
	}

	return nil
}

func (s *FollowerStore) Unfollow(ctx context.Context, followerID, userID int64) error {
	query := `
		DELETE FROM followers 
		WHERE user_id = $1 AND follower_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.Exec(ctx, query, userID, followerID)
	if err != nil {
		return fmt.Errorf("follower repo - unfollow: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFollowing
	}
	return nil
}
