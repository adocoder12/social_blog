package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/adocoder12/social_blog/internal/model"
	"github.com/adocoder12/social_blog/internal/store"
)

var ErrSelfFollow = errors.New("users cannot follow themselves")

type UserService struct {
	user    store.UserRepository
	follows store.FollowRepository
}

func NewUserService(user store.UserRepository, follows store.FollowRepository) *UserService {
	return &UserService{
		user:    user,
		follows: follows,
	}
}
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.user.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user service: get by id: %w", err)
	}
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	user, err := s.user.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("user service: Create user: %w", err)
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *model.User, patch model.UserPatch) (*model.User, error) {
	if patch.Username != nil {
		user.Username = *patch.Username
	}
	if patch.Email != nil {
		user.Email = *patch.Email

	}
	updated, err := s.user.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("user service: update user: %w", err)
	}
	return updated, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	if err := s.user.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("user service: delete user: %w", err)

	}
	return nil
}

func (s *UserService) Follow(ctx context.Context, followerID, userID int64) error {
	if followerID == userID {
		return ErrSelfFollow
	}
	if err := s.follows.Follow(ctx, followerID, userID); err != nil {
		return fmt.Errorf("user service: follow: %w", err)
	}
	return nil
}

func (s *UserService) Unfollow(ctx context.Context, followerID, userID int64) error {
	if err := s.follows.Unfollow(ctx, followerID, userID); err != nil {
		return fmt.Errorf("user service: unfollow: %w", err)
	}
	return nil
}
